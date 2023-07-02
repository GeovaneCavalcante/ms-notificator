from flask import Flask, request, jsonify
from flasgger import Swagger
import json
import requests

app = Flask(__name__)
Swagger(app)

notifications = []


@app.route('/', methods=['POST'])
def handle_notification():
    """
    Endpoint to handle notifications
    ---
    parameters:
      - in: body
        name: body
        required: true
        schema:
          id: Notifications
          required:
            - Type
          properties:
            Type:
              type: string
              description: The type of the message
            Message:
              type: string
              description: The message
            SubscribeURL:
              type: string
              description: The subscribe URL
    responses:
      200:
        description: OK
    """
    message = json.loads(request.data)
    print(message)

    if message['Type'] == 'SubscriptionConfirmation':
        print('Received notification: ' + message['Message'])
        print(message['SubscribeURL'].replace("localhost", "172.28.3.5"))
        requests.get(message['SubscribeURL'].replace(
            "localhost", "172.28.3.5"))

    elif message['Type'] == 'Notification':
        print('Received notification: ' + message['Message'])
        notifications.append(message)

    return '', 200


@app.route('/notifications', methods=['GET'])
def list_notifications():
    """
    Endpoint to list all notifications
    ---
    responses:
      200:
        description: A list of notifications
        schema:
          id: NotificationsList
          properties:
            notifications:
              type: array
              description: The notifications
              items:
                schema:
                  id: Notification
                  properties:
                    Type:
                      type: string
                      description: The type of the message
                    Message:
                      type: string
                      description: The message
                    SubscribeURL:
                      type: string
                      description: The subscribe URL
    """
    return jsonify(notifications), 200


if __name__ == '__main__':
    app.run(host="0.0.0.0", port=8080)
