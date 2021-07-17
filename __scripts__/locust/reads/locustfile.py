import logging
import uuid

import locust

from common import settings


logger = logging.getLogger()


class Reader(locust.HttpUser):
    host = settings.MEMKV_URL
    wait_time = locust.between(0, 0.1)

    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.cfg = settings.Configuration.get_instance()

    @locust.task
    def get(self) -> None:
        r = str(uuid.uuid4())
        template = '/keys/{key}'
        _ = self.client.get(
            template.format(key=r[:4]),
            name=template,
        )
