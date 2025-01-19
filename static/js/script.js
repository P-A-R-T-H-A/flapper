document.addEventListener("DOMContentLoaded", () => {
  const outputDiv = document.getElementById("output");
  const statusDiv = document.getElementById("status");
  const micButton = document.getElementById("micButton");

  let isListening = false; // Track active listening state
  const language = "en-US";
  const isEnglish = language === "en-US";

  // Audio feedback
  const startSound = new Audio("/static/sounds/start.mp3");
  const stopSound = new Audio("/static/sounds/stop.mp3");

  // Function to play sound safely
  function playSound(sound) {
    try {
      sound
        .play()
        .catch((error) => console.warn("Error playing sound:", error));
    } catch (error) {
      console.warn("Error playing sound:", error);
    }
  }

  // Check browser support for Web Speech API
  const SpeechRecognition =
    window.SpeechRecognition || window.webkitSpeechRecognition;
  if (!SpeechRecognition) {
    outputDiv.textContent = "Speech Recognition not supported in this browser.";
    return;
  }

  // Wake recognition (for en-US only)
  const wakeRecognition = isEnglish ? new SpeechRecognition() : null;
  if (wakeRecognition) {
    wakeRecognition.continuous = true;
    wakeRecognition.interimResults = true;
    wakeRecognition.lang = language;

    wakeRecognition.onresult = (event) => {
      const transcript = Array.from(event.results)
        .map((result) => result[0].transcript.toLowerCase().trim())
        .join(" ");

      if (transcript.includes("hello")) {
        wakeRecognition.stop();
        startActiveListening();
      }
    };

    wakeRecognition.onend = () => {
      if (!isListening) wakeRecognition.start(); // Restart only if not actively listening
    };
  }

  // Active listening
  const activeRecognition = new SpeechRecognition();
  activeRecognition.continuous = true;
  activeRecognition.interimResults = true;
  activeRecognition.lang = language;

  activeRecognition.onstart = () => {
    playSound(startSound);
    statusDiv.textContent = "Listening actively...";
    micButton.style.backgroundColor = "#FF6347"; // Red for active listening
  };

  activeRecognition.onresult = (event) => {
    const transcript = Array.from(event.results)
      .map((result) => result[0].transcript)
      .join("");
    outputDiv.textContent = transcript;
  };

  activeRecognition.onend = () => {
    playSound(stopSound);
    isListening = false;
    micButton.style.backgroundColor = "#4CAF50"; // Green for inactive
    statusDiv.textContent = isEnglish
      ? 'Say "Hello" to activate (en-US) or press the button'
      : "Press the button to start listening";
    if (isEnglish && !isListening) wakeRecognition.start(); // Restart wake word detection
  };

  // Start active listening
  function startActiveListening() {
    if (isListening) return; // Avoid multiple starts
    isListening = true;
    activeRecognition.start();
  }

  // Stop active listening
  function stopActiveListening() {
    if (!isListening) return;
    isListening = false;
    activeRecognition.stop();
  }

  // Toggle mic button functionality
  micButton.addEventListener("click", () => {
    if (isListening) {
      console.log("feature: call api");
      callApi(outputDiv.textContent);
      stopActiveListening();
    } else {
      if (isEnglish && wakeRecognition) wakeRecognition.stop(); // Stop wake detection
      startActiveListening();
    }
  });

  // Initialize wake word detection (en-US only)
  if (isEnglish && wakeRecognition) {
    wakeRecognition.start();
    statusDiv.textContent =
      'Say "Hello" to activate (en-US) or press the button';
  } else {
    statusDiv.textContent = "Press the button to start listening";
  }

  // call api http://localhost:8090/api/ai-agent/hotel-json-maker
  const callApi = async (voice) => {
    const response = await fetch(
      "http://localhost:8090/api/ai-agent/hotel-json-maker?voice=" + voice
    );
    const data = await response.json();
    console.log(data);
  };
});
