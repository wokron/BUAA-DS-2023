import argparse

from conn import ENSConnection


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--server", "-s", type=str, default="localhost")
    parser.add_argument("--port", "-p", type=int, default="4567")

    class process_key_calue(argparse.Action):
        def __call__(self, parser, namespace, values, option_string=None):
            if getattr(namespace, self.dest) is None:
                setattr(namespace, self.dest, list())

            for value in values:
                key, value = value.split("=")
                getattr(namespace, self.dest).append((key, value))

    subparsers = parser.add_subparsers()
    publish_parser = subparsers.add_parser("publish", help="publish events")
    publish_parser.add_argument(
        "--event", "-e", nargs="*", action=process_key_calue, help="event to publish"
    )
    publish_parser.set_defaults(name="publish")

    subscribe_parser = subparsers.add_parser(
        "subscribe", help="subscribe topics and wait for receiveing"
    )
    subscribe_parser.add_argument(
        "--topic", "-t", action="append", help="topic to subscribe"
    )
    subscribe_parser.set_defaults(name="subscribe")

    args = parser.parse_args()

    with ENSConnection(args.server, args.port) as conn:
        print(f"Connect to {args.server}:{args.port} success")

        if args.name == "publish":
            for topic, message in args.event:
                print(f'Publish event on topic "{topic}", message: "{message}"')
                conn.publish(topic, message)

            print(f"Publish all events")

        elif args.name == "subscribe":
            for topic in args.topic:
                print(f'Subscribe topic "{topic}"')
                conn.subscribe(topic)

            while True:
                msg = conn.recv()
                print(
                    f'Receive message on topic "{msg.topic}", message: "{msg.message}"'
                )
