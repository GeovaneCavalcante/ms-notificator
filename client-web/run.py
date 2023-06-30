from flask import Flask, request, jsonify
import json
import requests

app = Flask(__name__)

notifications = []

@app.route('/', methods=['POST'])
def handle_notification():
    message = json.loads(request.data)
    print(message)

    if message['Type'] == 'SubscriptionConfirmation':
        print('Received notification: ' + message['Message'])
        print(message['SubscribeURL'].replace("localhost", "172.20.0.5"))
        requests.get(message['SubscribeURL'].replace("localhost", "172.20.0.5"))
    
    elif message['Type'] == 'Notification':
        print('Received notification: ' + message['Message'])
        notifications.append(message)
    
    return '', 200

@app.route('/notifications', methods=['GET'])
def list_notifications():
    return jsonify(notifications), 200

if __name__ == '__main__':
    app.run(host="0.0.0.0", port=8080)