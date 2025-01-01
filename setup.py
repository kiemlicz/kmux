import setuptools

with open("README.md", "r") as fh:
    long_description = fh.read()
with open('requirements.txt') as f:
    requirements = f.read().splitlines()

setuptools.setup(
    name="kmux", 
    version="0.1.0",
    author="kiemlicz",
    author_email="stanislaw.dev@gmail.com",
    description="Kubeconfig manager",
    long_description=long_description,
    long_description_content_type="text/markdown",
    url="https://github.com/kiemlicz/kmux",
    scripts=['bin/km'],
    install_requires=requirements,
    packages=setuptools.find_packages(),
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
    ],
    python_requires='>=3.11',
)
