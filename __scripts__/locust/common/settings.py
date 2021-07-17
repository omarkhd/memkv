import copy
import hashlib
import json
import logging
import os
import threading
import time
from typing import Tuple

from common import network


logger = logging.getLogger()

SETTINGS_FILE = 'settings.json'
MEMKV_URL = os.getenv('LOCDEX_URL')

if MEMKV_URL is None:
    host = network.Host()
    MEMKV_URL = 'http://{}:4444'.format(host.gateway)


class FileWatcher(threading.Thread):
    _wait_seconds = 3

    def __init__(self, cfg):
        super().__init__()
        self.raw = None
        self.cfg = cfg

    def run(self):
        logger.info('Starting file watcher')
        first_run = True
        while True:
            if not first_run:
                time.sleep(self._wait_seconds)
            first_run = False

            try:
                f = open(SETTINGS_FILE, 'r')
                contents = f.read()
                f.close()
            except FileNotFoundError:
                logger.error('Settings file not found')
                continue

            try:
                loaded = json.loads(contents)
            except json.JSONDecodeError:
                logger.error('Invalid JSON file')
                continue

            if self.raw != contents:
                self.raw = contents
                algo = hashlib.md5()
                algo.update(self.raw.encode('utf-8'))
                self.cfg.update(loaded, algo.hexdigest())


class Configuration:
    _LOCK = threading.Lock()
    _INSTANCE = None

    def __init__(self):
        if not self._LOCK.locked():
            raise RuntimeError('global lock in wrong state')
        self._settings = None
        self._version = None
        self._lock = threading.Lock()
        self._updater = FileWatcher(self)
        self._updater.start()

    @classmethod
    def get_instance(cls):
        cls._LOCK.acquire()
        if cls._INSTANCE is None:
            cls._INSTANCE = cls()
        cls._LOCK.release()
        return cls._INSTANCE

    def update(self, settings: dict, version: str) -> bool:
        if not self.validate(settings):
            logger.error('Invalid settings object')
            return False
        self._lock.acquire()
        self._settings = settings
        self._version = version
        self._lock.release()
        logger.info('Configuration updated: %s', settings)
        return True

    def get(self) -> Tuple[dict, str]:
        self._lock.acquire()
        settings = copy.deepcopy(self._settings)
        version = self._version
        self._lock.release()
        return settings, version

    @staticmethod
    def validate(settings) -> bool:
        return type(settings) is dict
