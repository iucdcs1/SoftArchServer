const axios = require('axios');
const { performance } = require('perf_hooks');

const API_URL = 'http://79.174.95.21:8089';
const AUTH_TOKEN = 'NRn2vYTpx38iRyvJxAoQOuesJlcjEEiX';
const MESSAGE_THRESHOLD = 200; // 200 ms for message submission
const COUNT_THRESHOLD = 100;   // 100 ms for message count request
const LATENCY = 0;


async function measureTimeBehavior() {
  try {
    // Step 1: Measure time for sending a message
    const messageData = { text: 'Test message' };
    let startTime = performance.now();
    
    const messageResponse = await axios.post(`${API_URL}/messages`, messageData, {
      headers: {
        Authorization: `Bearer ${AUTH_TOKEN}`
      }
    });
    
    let endTime = performance.now();
    let messageResponseTime = endTime - startTime + LATENCY;

    console.log(`Message submission response time: ${messageResponseTime} ms`);

    // Step 2: Check if message response time meets the threshold
    let messageFitness = messageResponseTime <= MESSAGE_THRESHOLD ? 1 : (messageResponseTime / MESSAGE_THRESHOLD);
    if (messageResponseTime <= MESSAGE_THRESHOLD) {
      console.log(`Message submission was within the threshold of ${MESSAGE_THRESHOLD} ms.`);
    } else {
      console.log(`Message submission exceeded the threshold of ${MESSAGE_THRESHOLD} ms by ${messageResponseTime - MESSAGE_THRESHOLD} ms.`);
    }

    // Step 3: Measure time for requesting the message count
    startTime = performance.now();

    const countResponse = await axios.get(`${API_URL}/messages/count`, {
      headers: {
        Authorization: `Bearer ${AUTH_TOKEN}`
      }
    });

    endTime = performance.now();
    let countResponseTime = endTime - startTime + LATENCY;

    console.log(`Message count response time: ${countResponseTime} ms`);

    // Step 4: Check if message count response time meets the threshold
    let countFitness = countResponseTime <= COUNT_THRESHOLD ? 1 : (countResponseTime / COUNT_THRESHOLD);
    if (countResponseTime <= COUNT_THRESHOLD) {
      console.log(`Message count request was within the threshold of ${COUNT_THRESHOLD} ms.`);
    } else {
      console.log(`Message count request exceeded the threshold of ${COUNT_THRESHOLD} ms by ${countResponseTime - COUNT_THRESHOLD} ms.`);
    }

    // Step 5: Calculate the overall fitness score (lower is better)
    const overallFitness = (messageFitness + countFitness) / 2;

    console.log(`Overall fitness score: ${overallFitness}`);
    return overallFitness;
  } catch (error) {
    console.error('Error during time behavior measurement:', error);
  }
}

// Run the fitness function
measureTimeBehavior();
