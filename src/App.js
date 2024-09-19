import React, { useState, useEffect } from 'react';
import axios from 'axios';
import MessageList from './components/MessageList';
import MessageForm from './components/MessageForm';

const API_URL = 'http://79.174.95.21:8089';
const AUTH_TOKEN = 'NRn2vYTpx38iRyvJxAoQOuesJlcjEEiX';

const App = () => {
  const [messages, setMessages] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchMessages = async () => {
      try {
        const response = await axios.get(`${API_URL}/messages`, {
          headers: {
            Authorization: `Bearer ${AUTH_TOKEN}`
          }
        });
        setMessages(response.data);
        setLoading(false);
      } catch (error) {
        console.error('Error fetching messages:', error);
      }
    };

    fetchMessages();
  }, []);

  const addMessage = (newMessage) => {
    setMessages([...messages, newMessage]);
  };

  return (
    <div className="app">
      <h1>Anonymous Chat</h1>
      {loading ? (
        <p>Loading messages...</p>
      ) : (
        <MessageList messages={messages} />
      )}
      <MessageForm onMessageSend={addMessage} />
    </div>
  );
};

export default App;
