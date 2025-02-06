// Import the functions you need from the SDKs you need

import { initializeApp } from "firebase/app";
import { getAuth } from "firebase/auth";

// Your web app's Firebase configuration

const firebaseConfig = {
  apiKey: "AIzaSyAZ1MlWYNRjI1wNV8qGW7gZTUYpEmWqk8U",

  authDomain: "song-sleuths.firebaseapp.com",

  projectId: "song-sleuths",

  storageBucket: "song-sleuths.firebasestorage.app",

  messagingSenderId: "77751881239",

  appId: "1:77751881239:web:38335a29ccba33e5ac5652",
};

// Initialize Firebase

const app = initializeApp(firebaseConfig);

// Initialize Firebase Authentication and get a reference to the service
export const auth = getAuth(app);
