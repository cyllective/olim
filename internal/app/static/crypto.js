// Helper function to convert ArrayBuffer to base64
function arrayBufferToBase64(buffer) {
  const bytes = new Uint8Array(buffer);
  let binary = "";
  for (let i = 0; i < bytes.length; i++) {
    binary += String.fromCharCode(bytes[i]);
  }
  return btoa(binary);
}

// Helper function to convert base64 to ArrayBuffer
function base64ToArrayBuffer(base64) {
  const binary = atob(base64);
  const bytes = new Uint8Array(binary.length);
  for (let i = 0; i < binary.length; i++) {
    bytes[i] = binary.charCodeAt(i);
  }
  return bytes.buffer;
}

// Helper function to convert hex string to ArrayBuffer
function hexToArrayBuffer(hex) {
  const bytes = new Uint8Array(hex.length / 2);
  for (let i = 0; i < hex.length; i += 2) {
    bytes[i / 2] = parseInt(hex.substr(i, 2), 16);
  }
  return bytes.buffer;
}

// Helper function to convert ArrayBuffer to hex string
function arrayBufferToHex(buffer) {
  const bytes = new Uint8Array(buffer);
  return Array.from(bytes)
    .map((b) => b.toString(16).padStart(2, "0"))
    .join("");
}

// Generate a AES-GCM key
async function generateKey() {
  return await crypto.subtle.generateKey(
    {
      name: "AES-GCM",
      length: 256,
    },
    true, // extractable
    ["encrypt", "decrypt"],
  );
}

// Generate a IV (12 bytes)
function generateIV() {
  return crypto.getRandomValues(new Uint8Array(12));
}

// Encrypt data (text or file content)
async function encryptData(key, data) {
  const iv = generateIV();

  // Convert string to ArrayBuffer if needed
  const dataBuffer = typeof data === "string" ? new TextEncoder().encode(data) : data;

  const encryptedBuffer = await crypto.subtle.encrypt(
    {
      name: "AES-GCM",
      iv: iv,
    },
    key,
    dataBuffer,
  );

  return {
    encrypted: arrayBufferToBase64(encryptedBuffer),
    iv: arrayBufferToHex(iv),
  };
}

// Decrypt data
async function decryptData(keyHex, ivHex, encryptedBase64) {
  // Import the key
  const keyBuffer = hexToArrayBuffer(keyHex);
  const key = await crypto.subtle.importKey(
    "raw",
    keyBuffer,
    {
      name: "AES-GCM",
      length: 256,
    },
    false,
    ["decrypt"],
  );

  // Convert IV and encrypted data
  const iv = hexToArrayBuffer(ivHex);
  const encryptedBuffer = base64ToArrayBuffer(encryptedBase64);

  // Decrypt
  const decryptedBuffer = await crypto.subtle.decrypt(
    {
      name: "AES-GCM",
      iv: new Uint8Array(iv),
    },
    key,
    encryptedBuffer,
  );

  return decryptedBuffer;
}

// Export key as hex string
async function exportKeyAsHex(key) {
  const exportedKey = await crypto.subtle.exportKey("raw", key);
  return arrayBufferToHex(exportedKey);
}
