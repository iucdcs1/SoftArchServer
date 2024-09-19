# Project: Anonymous Chat Application

This project is a client-server application designed to facilitate anonymous messaging between users. Additionally, the system provides an endpoint to query the total number of messages sent. The project adheres to the assigned task and focuses on measuring three quality attributes: **Time Behavior**, **Recoverability**, and **Maintainability**.

## Repository Overview
The repository consists of three branches:

- **backend**: Contains the server-side code that manages the chat functionality and message count endpoint.
- **frontend**: Contains the client-side code for the chat application interface.
- **main**: Includes the **fitness functions** designed to evaluate the system based on the quality attributes outlined in the task.

### Branch Navigation

1. **backend branch**:
   - Server code responsible for handling message transmission and retrieval.
   - Exposes three endpoint: to get messages, sent new or to return the total number of messages sent.

2. **frontend branch**:
   - Client code for user interaction.
   - Provides the interface for sending and displaying messages in the prototype of the chatroom.

3. **main branch**:
   - Implements fitness functions for testing the quality attributes:
     - **Time Behavior**: Measures response time for message submissions and the `/messages/count` endpoint.
     - **Recoverability**: Tests the system’s ability to recover from crashes without message loss.
     - **Maintainability**: Evaluates the system's modularity and ease of refactoring without breaking functionality.

## Task Brief Overview

This project was developed in response to the following task:

- **Anonymous Messaging**: Users can send messages anonymously, which are displayed alongside the date and time sent.
- **Message Count Endpoint**: A GET request to `/messages/count` returns the current number of messages.
- **Quality Attribute Scenarios**: The project focuses on evaluating **Time Behavior**, **Recoverability**, and **Maintainability** through defined scenarios and fitness functions.

### Quality Attribute Scenarios
1. **Time Behavior**: The system must respond to message submissions within 200ms and return the message count within 100ms under normal conditions.
2. **Recoverability**: The server should recover within 30 seconds following a crash, without losing messages, and reconnect clients automatically.
3. **Maintainability**: Refactoring the codebase should maintain functionality and pass automated tests with no more than two additional test cases.

---

Feel free to explore each branch, review the code of the client- or server- implementation, and check the fitness function results in the `main` branch to evaluate the system's performance.
