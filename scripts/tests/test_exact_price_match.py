import time
import requests
from . import API_URL, newUUID


def run():
    print("Testing Exact Price Match Flow...")
    symbol = newUUID()
    price = 700.0

    # Create Brokers
    b1_resp = requests.post(f"{API_URL}/brokers", json={"name": "Broker 1"})
    b1_id = b1_resp.json().get("id")
    b2_resp = requests.post(f"{API_URL}/brokers", json={"name": "Broker 2"})
    b2_id = b2_resp.json().get("id")

    ask_payload = {
        "broker_id": b1_id,
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
        "broker_id": b2_id,
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
