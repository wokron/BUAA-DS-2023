MAX_TOPIC_LENGTH = 50
MAX_MESSAGE_LENGTH = 50


class ENSMsg:
    def __init__(self, type: int, topic: str = "", message: str = ""):
        self.type = type
        self.topic = topic
        self.message = message

    def encode(self) -> bytes:
        type_data = self.type.to_bytes(1, byteorder="big")

        topic_data = bytes(self.topic, "utf-8")
        topic_data = self.resize_to_length(topic_data, MAX_TOPIC_LENGTH)
        assert len(topic_data) == MAX_TOPIC_LENGTH

        message_data = bytes(self.message, "utf-8")
        message_data = self.resize_to_length(message_data, MAX_MESSAGE_LENGTH)
        assert len(message_data) == MAX_MESSAGE_LENGTH

        return type_data + topic_data + message_data

    @staticmethod
    def resize_to_length(data: bytes, length: int) -> bytes:
        data = data[:length] + bytes([0] * max(0, length - len(data)))
        return data

    @classmethod
    def decode(cls, data: bytes):
        type_data = data[0]
        topic_data = data[1 : 1 + MAX_TOPIC_LENGTH]
        message_data = data[1 + MAX_TOPIC_LENGTH :]

        return cls(type_data, topic_data.decode("utf-8"), message_data.decode("utf-8"))
