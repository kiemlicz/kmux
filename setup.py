import setuptools

with open("README.md", "r") as fh:
    long_description = fh.read()

setuptools.setup(
    name="kmux", 
    version="0.0.1",
    author="kiemlicz",
    author_email="stanislaw.dev@gmail.com",
    description="Kubeconfig Multiplexer",
    long_description=long_description,
    long_description_content_type="text/markdown",
    url="https://github.com/kiemlicz/kmux",
    scripts=['bin/km'],
    install_requires=[
        "argparse~=1.4.0",
        "PyYAML~=5.3.1",
        "libtmux~=0.8.3",
        "setuptools~=50.3.0",
        "google-cloud-container~=2.1.0",
        "google-api-python-client~=1.12.3",
        "google-auth-oauthlib~=0.4.1",
        "google-auth-httplib2~=0.0.4"
    ],
    packages=setuptools.find_packages(),
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
    ],
    python_requires='>=3.6',
)
