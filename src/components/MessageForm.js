import React, { useState } from 'react';
import axios from 'axios';

const API_URL = 'http://79.174.95.21:8089';
const AUTH_TOKEN = 'NRn2vYTpx38iRyvJxAoQOuesJlcjEEiX'; 

const MessageForm = ({ onMessageSend }) => {
  const [text, setText] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (text.trim() === '') return;

    try {
      const response = await axios.post(`${API_URL}/messages`, { text }, {
        headers: {
          Authorization: `Bearer ${AUTH_TOKEN}`
        }
      });

      onMessageSend(response.data);
      setText('');
    } catch (error) {
      console.error('Error sending message:', error);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="message-form">
      <input
        type="text"
        value={text}
        onChange={(e) => setText(e.target.value)}
        placeholder="Enter your message"
      />
      <button type="submit">Send</button>
    </form>
  );
};

export default MessageForm;
