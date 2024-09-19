const axios = require('axios');
const { exec } = require('child_process');
const { performance } = require('perf_hooks');

const API_URL = 'http://79.174.95.21:8089';
const AUTH_TOKEN = 'NRn2vYTpx38iRyvJxAoQOuesJlcjEEiX';
const RECOVERY_THRESHOLD = 30000; // 30 seconds
const MESSAGE_INTERVAL = 100; // 5 seconds for sending messages
let COUNT = 0

const MESSAGE_QUEUE = [];
let recoveryStartTime;
let messageIntervalId = null;  // Store the interval ID here

async function sendMessage(message) {
  try {
    await axios.post(`${API_URL}/messages`, message, {
      headers: {
        Authorization: `Bearer ${AUTH_TOKEN}`
      }
    });
    console.log('Message sent successfully:', message);
  } catch (error) {
    console.error('Failed to send message:', message, error.message); 
  }
  MESSAGE_QUEUE.push(message);
}

function startSendingMessages() {
  messageIntervalId = setInterval(() => {
    const message = { text: `Message test ${COUNT}` };  // Create a new message with the updated COUNT
    COUNT = COUNT + 1;
    sendMessage(message);
  }, MESSAGE_INTERVAL);
}

function stopSendingMessages() {
  if (messageIntervalId !== null) {
    clearInterval(messageIntervalId);
    console.log('Message sending stopped after server recovery.');
  }
}

function checkServerAvailability() {
  return new Promise((resolve, reject) => {
    exec(`curl -s -o /dev/null -w "%{http_code}" -H "Authorization: Bearer ${AUTH_TOKEN}" ${API_URL}/messages`, (error, stdout) => {
      if (error) {
        reject(error);
        return;
      }

      // Check if the status code is 200
      const statusCode = parseInt(stdout.trim(), 10);
      if (statusCode === 200) {
        resolve(true);
      } else {
        resolve(false);
      }
    });
  });
}

function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

async function testRecoverability() {
  try {
    // Simulate server crash (manually or programmatically)
    console.log('Simulating server crash... Please restart the server manually.');

    // Wait until server goes down (server crash)
    let serverCrashed = false;
    console.log('Waiting for the server to crash...');
    while (!serverCrashed) {
      try {
        const serverAvailable = await checkServerAvailability();
        if (!serverAvailable) {
          serverCrashed = true;
          console.log('Server has crashed.');
        } else {
          console.log('Server is still running...');
          await sleep(1000); // Wait for 1 second before retrying
        }
      } catch (error) {
        serverCrashed = true;
        console.log('Server has crashed (or became unavailable).');
      }
    }

    // Start sending messages after the crash
    console.log('Starting to send messages after server crash...');
    startSendingMessages();

    // After the crash, begin recovery process
    const startTime = performance.now();
    let serverRecovered = false;

    while (!serverRecovered && (performance.now() - startTime < RECOVERY_THRESHOLD)) {
      try {
        serverRecovered = await checkServerAvailability();
        if (serverRecovered) {
          console.log('Server has recovered.');
        } else {
          console.log('Waiting for server to recover...');
          await sleep(1000); // Wait for 1 second before retrying
        }
      } catch (error) {
        console.log('Waiting for server to recover...');
        await sleep(1000); // Wait for 1 second before retrying
      }
    }

    if (performance.now() - startTime >= RECOVERY_THRESHOLD) {
      console.log('Server not recovered after 30s.');
    } else {
      const endTime = performance.now();
      console.log('Server recovered successfully.');

      // Stop sending messages after recovery
      stopSendingMessages();

      // Ensure no messages were lost
      const messagesResponse = await axios.get(`${API_URL}/messages`, {
        headers: {
          Authorization: `Bearer ${AUTH_TOKEN}`
        }
      });

      const messages = messagesResponse.data;
      const messageTexts = messages.map(msg => msg.text);

      // Check if all messages in MESSAGE_QUEUE exist
      let allMessagesExist = true;

      for (const queuedMessage of MESSAGE_QUEUE) {
        if (!messageTexts.includes(queuedMessage.text)) {
          allMessagesExist = false;
          console.log(`Message was lost during the crash: ${queuedMessage.text}`);
        }
      }

      if (allMessagesExist) {
        console.log('No messages were lost during the crash.');
      } else {
        console.log('Some messages were lost during the crash.');
      }

      // Check if recovery time is within the threshold
      const recoveryTime = endTime - startTime;
      if (recoveryTime <= RECOVERY_THRESHOLD) {
        console.log(`Recovery time (${recoveryTime}) is within the acceptable threshold of ${RECOVERY_THRESHOLD} ms.`);
      } else {
        console.log(`Recovery time exceeded the threshold by ${recoveryTime - RECOVERY_THRESHOLD} ms.`);
      }
    }

  } catch (error) {
    console.error('Error during recoverability test:', error.message);
  }
}

// Run the recoverability test
testRecoverability();
