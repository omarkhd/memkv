import logging
import random
import uuid

import locust

from common import settings


logger = logging.getLogger()


class Writer(locust.HttpUser):
    host = settings.MEMKV_URL
    wait_time = locust.between(0, 1)

    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.cfg = settings.Configuration.get_instance()
        self.item_count = 0  # Used to stop writing
        self.item_cache = None  # Used to cache one key
        self.limit = -1

    @locust.task(8)
    def put(self) -> None:
        if self.item_count > self.limit:
            return

        r = str(uuid.uuid4())
        template = '/keys/{key}'
        _ = self.client.post(
            template.format(key=r[:4]),
            data=r, name=template,
        )

    @locust.task(1)
    def keys(self) -> None:
        # Updating item count and cache
        r = self.client.get('/keys')
        if r.status_code == 200:
            keys = r.json()
            if len(keys) > 0:
                self.item_count = len(keys)
                self.item_cache = random.choice(keys)

        # Updating limit from configuration
        cfg, _ = self.cfg.get()
        limit = cfg.get('limit')
        if type(limit) is int:
            self.limit = limit

    @locust.task(1)
    def delete(self) -> None:
        template = '/keys/{key}'
        _ = self.client.delete(
            template.format(key=self.item_cache),
            name=template,
        )
