import time
import requests
from . import API_URL, newUUID

def run():
    print("Testing Multiple Matches Flow...")
    symbol = newUUID()
    ask_price = 200.0
    bid_price = 210.0

    # 1. Create two smaller Ask Orders
    asks = []
    for i in range(2):
        payload = {
            "broker_id": f"broker_ask{i}",
            "owner_doc": f"doc_ask{i}",
            "type": "ASK",
            "symbol": symbol,
            "price": ask_price,
            "quantity": 5,
            "valid_until": time.strftime("%Y-%m-%d")
        }
        resp = requests.post(f"{API_URL}/orders", json=payload)
        assert resp.status_code == 201
        asks.append(resp.json().get("id"))

    # 2. Create one larger Bid Order that covers both
    bid_payload = {
        "broker_id": "broker_bid",
        "owner_doc": "doc_bid",
        "type": "BID",
        "symbol": symbol,
        "price": bid_price,
        "quantity": 10,
        "valid_until": time.strftime("%Y-%m-%d")
    }
    resp = requests.post(f"{API_URL}/orders", json=bid_payload)
    assert resp.status_code == 201
    bid_id = resp.json().get("id")

    time.sleep(2)

    # 3. Verify all filled
    for ask_id in asks:
        order = requests.get(f"{API_URL}/orders/{ask_id}").json()
        print(f"Ask {ask_id} status: {order.get('status')}")
        assert order.get("status") == "FILLED", f"Ask order {ask_id} not filled: {order}"
    
    bid_order = requests.get(f"{API_URL}/orders/{bid_id}").json()
    print(f"Bid {bid_id} status: {bid_order.get('status')}")
    assert bid_order.get("status") == "FILLED", f"Bid order not filled: {bid_order}"
    print("Multiple matches successful")
