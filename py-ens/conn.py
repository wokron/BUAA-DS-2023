import socket

from msg import (
    MAX_MESSAGE_LENGTH,
    MAX_TOPIC_LENGTH,
    ENSMsg,
)

SUBSCRIBE = 0
UNSUBSCRIBE = 1
PUBLISH = 2
UPDATE = 3


class ENSConnection:
    def __init__(self, host, port):
        self.stream = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.stream.connect((host, port))

    def __enter__(self):
        return self

    def __exit__(self, *args):
        self.close()

    def close(self):
        self.stream.close()

    def send(self, msg: ENSMsg):
        self.stream.send(msg.encode())

    def recv(self) -> ENSMsg:
        data = self.stream.recv(1 + MAX_TOPIC_LENGTH + MAX_MESSAGE_LENGTH)
        return ENSMsg.decode(data)

    def publish(self, topic: str, message: str):
        msg = ENSMsg(PUBLISH, topic, message)
        self.send(msg)

    def subscribe(self, topic: str):
        msg = ENSMsg(SUBSCRIBE, topic)
        self.send(msg)

    def unsubscribe(self, topic: str):
        msg = ENSMsg(UNSUBSCRIBE, topic)
        self.send(msg)
