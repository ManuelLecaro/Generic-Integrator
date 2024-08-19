from flask import Flask, jsonify, request
import uuid

app = Flask(__name__)

# Dictionary to simulate a database of payments
payments_db = {}

@app.route('/pay', methods=['POST'])
def pay():
    transaction_id = str(uuid.uuid4())

    payments_db[transaction_id] = {
        "transaction_id": transaction_id,
        "message": "Payment processed successfully."
    }
    
    response = {
        "transaction_id": transaction_id,
        "message": "Payment processed successfully."
    }
    
    return jsonify(response), 200

@app.route('/refund', methods=['POST'])
def refund():
    data = request.get_json()
    transaction_id = data.get('transaction_id')
    
    response = {
        "transaction_id": transaction_id,
        "message": "Refund processed successfully."
    }
    
    return jsonify(response), 200

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
