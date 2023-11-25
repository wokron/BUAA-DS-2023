use std::net::TcpStream;

use enslib::{client::subscribe, conn::recv_ens_msg};

mod enslib;

fn main() {
    let mut stream = TcpStream::connect("localhost:4567").unwrap();

    let topics = vec![
        String::from("topic1"),
        String::from("topic2"),
        String::from("topic3"),
    ];

    for topic in topics {
        if let Err(err) = subscribe(&mut stream, &topic) {
            println!("Fail to subscribe topic \"{}\", {}", topic, err);
        }
    }

    loop {
        if let Some(msg) = recv_ens_msg(&mut stream) {
            println!(
                "Received Event, topic: \"{}\", message: \"{}\".",
                msg.topic, msg.message
            );
        } else {
            println!("error when receiving")
        }
    }
}
