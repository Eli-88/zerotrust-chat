import os
import subprocess

FILE_DIR = os.path.abspath(os.path.dirname(__file__))
ROOT_DIR = os.path.abspath(os.path.join(FILE_DIR, '../'))

if __name__ == '__main__':
    os.environ['GOPATH'] = ROOT_DIR
    os.environ['PATH'] = "{}:{}/bin".format(
        os.environ['PATH'], os.environ['GOPATH'])
    os.chdir(path=ROOT_DIR + '/src')
    subprocess.call('go generate ./...', shell=True)
