use std::{
    io::{Read, Write},
    net::TcpStream,
};

use super::msg::{ENSMsg, ENS_MGS_SIZE};
use std::io::Result;

pub fn send_ens_msg(stream: &mut TcpStream, msg: &ENSMsg) -> Result<()> {
    let data = msg.encode();

    stream.write_all(data.as_slice())?;

    return Ok(());
}

pub fn recv_ens_msg(stream: &mut TcpStream) -> Option<ENSMsg> {
    let mut recv_data: [u8; 101] = [0u8; ENS_MGS_SIZE];

    if let Err(_) = stream.read_exact(&mut recv_data) {
        return None;
    }

    return ENSMsg::decode(&recv_data.to_vec());
}
