import logging
import os
import pickle
from urllib.parse import urlparse, parse_qs
from pathlib import Path
from google_auth_oauthlib.flow import InstalledAppFlow
from googleapiclient.discovery import build
from google.auth.transport.requests import Request


log = logging.getLogger()
log.setLevel(logging.INFO)

SCOPES = [
    'https://www.googleapis.com/auth/drive.metadata.readonly',
    'https://www.googleapis.com/auth/drive.readonly',
]


class GDriveClient(object):
    def __init__(self, credentials_file):
        credentials_file = os.path.expanduser(credentials_file)
        cred_dir = Path(credentials_file).parent
        cached_credentials = os.path.join(cred_dir, 'token.pickle')
        credentials = None
        if os.path.exists(cached_credentials):
            with open(cached_credentials, 'rb') as token:
                credentials = pickle.load(token)
        if not credentials or not credentials.valid:
            if credentials and credentials.expired and credentials.refresh_token:
                credentials.refresh(Request())
            else:
                flow = InstalledAppFlow.from_client_secrets_file(os.path.expanduser(credentials_file), SCOPES)
                credentials = flow.run_local_server(port=0)
            with open(cached_credentials, 'wb') as token:
                pickle.dump(credentials, token)

        self.service = build('drive', 'v3', credentials=credentials, cache_discovery=False)

    def export(self, file_id, mime_type=None):
        '''
        Exports Google Doc to mime_type
        :param file_id: Google Drive's file ID
        :param mime_type: user-picked mime type
        :return: File contents in requested mime type
        '''
        exported = self.service.files().export(fileId=file_id, mimeType=mime_type).execute()
        return exported  # todo decode?

    def download(self, file_id):
        '''
        Download file as-is
        :param file_id: the Google file ID
        :return: file contents without any conversion
        '''
        exported = self.service.files().get_media(fileId=file_id).execute()
        return exported

    def get_file_contents(self, path_qs):
        log.debug("finding file: {}".format(path_qs))
        query_parsed = urlparse(path_qs)
        query_dict = self._query_string_to_dict(query_parsed.query)
        mime_type = query_dict['mime_type'] if "mime_type" in query_dict else None
        file_meta = self.get_file_meta(path_qs)
        return self.get_file(file_meta, mime_type)

    def get_file(self, file_meta, mime_type=None):
        if file_meta['mimeType'] == 'application/vnd.google-apps.document':
            log.debug("Exporting file: {}".format(file_meta))
            return self.export(file_meta['id'], mime_type if mime_type else file_meta['mimeType'])
        else:
            log.debug("Downloading file: {}".format(file_meta))
            return self.download(file_meta['id'])

    def get_file_meta(self, path_qs):
        '''
        Asserts that path exists on the google drive

        :return: full file_meta of file/folder traversed to (the last one)
        '''
        path_segment_list = self._path_to_list(path_qs)
        log.debug("gdrive segment list: {}".format(path_segment_list))

        def go(parent_meta, idx):
            if idx >= len(path_segment_list):
                return parent_meta
            next_name = path_segment_list[idx]
            file_list = self._list_children(parent_meta)
            r = [e for e in file_list if e['name'] == next_name]
            if len(r) > 0:
                # don't care if name occurred in other pages or already multiple times
                return go(r[0], idx + 1)
            raise ValueError('Unable to find name: {}, under directory with meta: {}'.format(next_name, parent_meta))

        if not path_segment_list:
            return {'id': 'root', 'mimeType': ''}
        else:
            return go({'id': 'root'}, 0)

    def _list_children(self, parent_meta):
        def query(extra_params={}):
            r = self.service.files().list(q="'{}' in parents and trashed = false".format(parent_meta['id']),
                                          **extra_params).execute()
            self._assert_incomplete_search(r)
            return r

        json_response = query()
        ret_list = json_response['files']
        while 'nextPageToken' in json_response:
            log.debug("Fetching next page of files under: {}".format(parent_meta))
            json_response = query({'pageToken': json_response['nextPageToken']})
            ret_list.extend(json_response['files'])
        return ret_list

    def _path_to_list(self, path):
        source = urlparse(path)
        p = source.netloc + source.path
        return p.strip(os.sep).split(os.sep)

    def _query_string_to_dict(self, qs):
        query_dict = parse_qs(qs)

        # unwrap values as by default arguments are parsed to lists { 'k': ['v']}
        for k, v in query_dict.items():
            if isinstance(v, list) and len(v) == 1:
                query_dict[k] = v[0]

        return query_dict

    def _assert_incomplete_search(self, json_response):
        if json_response['incompleteSearch']:
            raise RuntimeError('google drive query ended due to incompleteSearch')
