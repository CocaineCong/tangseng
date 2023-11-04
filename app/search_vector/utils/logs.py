"""
the log file 
"""
import os
import re
import datetime
import logging
import sys

try:
    import codecs
except ImportError:
    codecs = None


class MultiprocessHandler(logging.FileHandler):
    """
    log hander objective
    """

    def __init__(self,
                 filename,
                 when='D',
                 backup_count=0,
                 encoding=None,
                 delay=False):
        self.prefix = filename
        self.backup_count = backup_count
        self.when = when.upper()
        self.ext_math = r"^\d{4}-\d{2}-\d{2}"

        self.when_dict = {
            'S': "%Y-%m-%d-%H-%M-%S",
            'M': "%Y-%m-%d-%H-%M",
            'H': "%Y-%m-%d-%H",
            'D': "%Y-%m-%d"
        }

        self.suffix = self.when_dict.get(when)
        if not self.suffix:
            print('The specified date interval unit is invalid: ', self.when)
            sys.exit(1)

        self.filefmt = os.path.join('.', "logs",
                                    f"{self.prefix}-{self.suffix}.log")
        self.file_path = datetime.datetime.now().strftime(self.filefmt)
        _dir = os.path.dirname(self.filefmt)
        try:
            if not os.path.exists(_dir):
                os.makedirs(_dir)
        except Exception as e:
            print('Failed to create log file: ', e)
            print("log_path:" + self.file_path)
            sys.exit(1)

        if codecs is None:
            encoding = None

        logging.FileHandler.__init__(self, self.file_path, 'a+', encoding,
                                     delay)

    def should_change_file_to_write(self):
        """
        judge if need to change file to store log message
        """
        _file_path = datetime.datetime.now().strftime(self.filefmt)
        if _file_path != self.file_path:
            self.file_path = _file_path
            return True
        return False

    def do_change_file(self):
        """
        change file to store log message
        """
        self.base_filename = os.path.abspath(self.file_path)
        if self.stream:
            self.stream.close()
            self.stream = None

        if not self.delay:
            self.stream = self._open()
        if self.backup_count > 0:
            for s in self.get_files_to_delete():
                os.remove(s)

    def get_files_to_delete(self):
        dir_name, _ = os.path.split(self.baseFilename)
        file_names = os.listdir(dir_name)
        result = []
        prefix = self.prefix + '-'
        for file_name in file_names:
            if file_name[:len(prefix)] == prefix:
                suffix = file_name[len(prefix):-4]
                if re.compile(self.ext_math).match(suffix):
                    result.append(os.path.join(dir_name, file_name))
        result.sort()

        if len(result) < self.backup_count:
            result = []
        else:
            result = result[:len(result) - self.backup_count]
        return result

    def emit(self, record):
        try:
            if self.should_change_file_to_write():
                self.do_change_file()
            logging.FileHandler.emit(self, record)
        except Exception as e:
            LOGGER.error(f"emit error {e}")
            self.handleError(record)


def write_log():
    """
    store log
    """
    logger = logging.getLogger()
    logger.setLevel(logging.DEBUG)
    fmt = logging.Formatter('%(asctime)s | %(levelname)s | %(filename)s '
                            '| %(funcName)s | %(lineno)s | %(message)s')

    stream_handler = logging.StreamHandler(sys.stdout)
    stream_handler.setLevel(logging.INFO)
    stream_handler.setFormatter(fmt)

    log_name = "milvus"
    file_handler = MultiprocessHandler(log_name, when='D', backup_count=2)
    file_handler.setLevel(logging.DEBUG)
    file_handler.setFormatter(fmt)
    file_handler.do_change_file()

    logger.addHandler(stream_handler)
    logger.addHandler(file_handler)

    return logger


LOGGER = write_log()
