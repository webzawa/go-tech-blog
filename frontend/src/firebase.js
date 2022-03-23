// import firebase from 'firebase/app';
import firebase from 'firebase/compat/app';
import 'firebase/compat/auth';

const firebaseConfig = {
  apiKey: "AIzaSyD5iVXWumP_Yxkzlu66Zd3uu0VVacL3kCw",
  authDomain: "seisou-75737.firebaseapp.com",
  projectId: "seisou-75737",
  storageBucket: "seisou-75737.appspot.com",
  messagingSenderId: "312272401234",
  appId: "1:312272401234:web:06088f040cdf53f77ec08b",
  measurementId: "G-0D935W0ZDT"
};

firebase.initializeApp(firebaseConfig);

export const auth = firebase.auth();
