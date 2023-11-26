use std::net::TcpStream;

use super::{
    conn::send_ens_msg,
    msg::{ENSMsg, ENSMsgType},
};
use std::io::Result;

pub fn publish(stream: &mut TcpStream, topic: &String, message: &String) -> Result<()> {
    let msg = ENSMsg {
        msg_type: ENSMsgType::PUBLISH as u8,
        topic: topic.to_string(),
        message: message.to_string(),
    };

    return send_ens_msg(stream, &msg);
}

pub fn subscribe(stream: &mut TcpStream, topic: &String) -> Result<()> {
    let msg = ENSMsg {
        msg_type: ENSMsgType::SUBSCRIBE as u8,
        topic: topic.to_string(),
        message: String::new(),
    };

    return send_ens_msg(stream, &msg);
}

pub fn unsubscribe(stream: &mut TcpStream, topic: &String) -> Result<()> {
    let msg = ENSMsg {
        msg_type: ENSMsgType::UNSUBSCRIBE as u8,
        topic: topic.to_string(),
        message: String::new(),
    };

    return send_ens_msg(stream, &msg);
}
