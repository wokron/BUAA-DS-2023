use std::net::TcpStream;

use args::Args;
use clap::Parser;
use enslib::{
    client::{publish, subscribe},
    conn::recv_ens_msg,
};
mod args;

#[allow(dead_code)]
mod enslib;

fn main() {
    let args = Args::parse();

    let addr: String = args.server + ":" + &args.port.to_string();

    if let Ok(mut stream) = TcpStream::connect(addr) {
        match args.command {
            args::Commands::Publish { event } => {
                handle_publish(&mut stream, event);
            }
            args::Commands::Subscribe { topic } => {
                handle_subscribe(&mut stream, topic);
            }
        }
    } else {
        println!("Couldn't connect to server");
    }
}

fn handle_subscribe(stream: &mut TcpStream, topics: Vec<String>) {
    for topic in &topics {
        if let Err(err) = subscribe(stream, topic) {
            println!("Fail to subscribe topic \"{}\", {}", topic, err);
        }
    }

    loop {
        if let Ok(msg) = recv_ens_msg(stream) {
            println!(
                "Receive event, topic: \"{}\", message: \"{}\".",
                msg.topic, msg.message
            );
        } else {
            println!("error when receiving")
        }
    }
}

fn handle_publish(stream: &mut TcpStream, events: Vec<(String, String)>) {
    for event in &events {
        let topic = &event.0;
        let message = &event.1;

        if let Err(err) = publish(stream, topic, message) {
            println!(
                "Fail to publish event on topic \"{}\", message: \"{}\". {}.",
                topic, message, err
            );
        }
    }
    println!("Publishing events finished.")
}
