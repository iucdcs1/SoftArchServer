import React from 'react';

const MessageList = ({ messages }) => {
  return (
    <div className="message-list">
      {messages.map((message, index) => (
        <div key={index} className="message-item">
          <p><strong>{message.created_at}</strong>: {message.text}</p>
        </div>
      ))}
    </div>
  );
};

export default MessageList;
