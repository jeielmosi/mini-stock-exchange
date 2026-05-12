import time
import requests
from . import API_URL, newUUID


def run():
    print("Testing Exact Price Match Flow...")
    symbol = newUUID()
    price = 700.0

    ask_payload = {
        "broker_id": "b1",
        "owner_doc": "d1",
        "type": "ASK",
        "symbol": symbol,
        "price": price,
        "quantity": 10,
        "valid_until": time.strftime("%Y-%m-%d")
    }
    resp = requests.post(f"{API_URL}/orders", json=ask_payload)
    ask_id = resp.json().get("id")

    bid_payload = {
        "broker_id": "b2",
        "owner_doc": "d2",
        "type": "BID",
        "symbol": symbol,
        "price": price,
        "quantity": 10,
        "valid_until": time.strftime("%Y-%m-%d")
    }
    resp = requests.post(f"{API_URL}/orders", json=bid_payload)
    bid_id = resp.json().get("id")

    time.sleep(1)

    assert requests.get(f"{API_URL}/orders/{ask_id}").json().get("status") == "FILLED"
    assert requests.get(f"{API_URL}/orders/{bid_id}").json().get("status") == "FILLED"
    print("Exact price match successful")
