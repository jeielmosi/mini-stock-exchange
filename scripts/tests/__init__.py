import base64
import uuid
from utils import API_URL

def newUUID():
    return base64.urlsafe_b64encode(uuid.uuid4().bytes).rstrip(b'=').decode('ascii')[:20]
