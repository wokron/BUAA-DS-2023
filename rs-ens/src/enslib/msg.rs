use anyhow::{Result, Ok};
use std::{cmp, mem};

const MAX_TOPIC_LENGTH: usize = 50;
const MAX_MESSAGE_LENGTH: usize = 50;

pub const ENS_MGS_SIZE: usize = mem::size_of::<ENSMsgData>();

#[repr(u8)]
pub enum ENSMsgType {
    SUBSCRIBE,
    UNSUBSCRIBE,
    PUBLISH,
    UPDATE,
}

#[repr(C)]
struct ENSMsgData {
    msg_type: u8,
    topic: [u8; MAX_TOPIC_LENGTH],
    message: [u8; MAX_MESSAGE_LENGTH],
}

impl ENSMsgData {
    fn encode(&self) -> Vec<u8> {
        let data_slice: &[u8] = unsafe {
            let ptr = self as *const ENSMsgData as *const u8;
            std::slice::from_raw_parts(ptr, mem::size_of::<ENSMsgData>())
        };
        let mut data: Vec<u8> = Vec::new();
        data.extend_from_slice(data_slice);
        return data;
    }

    fn decode(data: &Vec<u8>) -> ENSMsgData {
        let slice_data = &data[..];
        let msg: ENSMsgData = unsafe { std::ptr::read(slice_data.as_ptr() as *const _) };
        msg
    }

    fn to_msg(&self) -> Result<ENSMsg> {
        let topic_str = String::from_utf8(self.topic.to_vec())?;
        let message_str = String::from_utf8(self.message.to_vec())?;

        let msg = ENSMsg {
            msg_type: self.msg_type,
            topic: topic_str,
            message: message_str,
        };

        Ok(msg)
    }
}

pub struct ENSMsg {
    pub msg_type: u8,
    pub topic: String,
    pub message: String,
}

impl ENSMsg {
    fn to_msg_data(&self) -> ENSMsgData {
        let mut msg_data = ENSMsgData {
            msg_type: self.msg_type,
            topic: [0u8; 50],
            message: [0u8; 50],
        };
        let src_message = self.message.as_bytes();
        let src_topic = self.topic.as_bytes();
        let min_message_len = cmp::min(src_message.len(), MAX_MESSAGE_LENGTH);
        let min_topic_len = cmp::min(src_topic.len(), MAX_TOPIC_LENGTH);

        msg_data.message[..min_message_len].copy_from_slice(&src_message[..min_message_len]);
        msg_data.topic[..min_topic_len].copy_from_slice(&src_topic[..min_topic_len]);
        msg_data
    }

    pub fn encode(&self) -> Vec<u8> {
        let msg_data = self.to_msg_data();
        let mut data: Vec<u8> = Vec::new();

        data.extend(msg_data.encode());
        data
    }

    pub fn decode(data: &Vec<u8>) -> Result<ENSMsg> {
        let msg_data = ENSMsgData::decode(data);
        msg_data.to_msg()
    }
}
