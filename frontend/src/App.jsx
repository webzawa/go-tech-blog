import React, { Component } from 'react';
import firebase from './firebase';
import SignInScreen from './components/SignInScreen';
import axios from 'axios';
// axios.defaults.baseURL = 'http://0.0.0.0:8080';
// axios.defaults.headers.post['Content-Type'] = 'application/json;charset=utf-8';
// axios.defaults.headers.post['Access-Control-Allow-Origin'] = '*';

class App extends Component {
  state = {
    loading: true,
    user: null
  };

  componentDidMount() {
    firebase.auth().onAuthStateChanged(user => {
      this.setState({
        loading: false,
        user: user
      });

      // axios.post(`http://0.0.0.0:8080/api/users`, {
      //   uid: "this.state.user.displayName",
      //   name: "this.state.user.uid"
      // }, {
      //   headers: {
      //     // "Access-Control-Allow-Origin": true,
      //     // 'Content-Type': 'application/x-www-form-urlencoded'
      //     "Content-Type": "application/json"
      //   }
      // })
      // .then(res => {
      //   console.log("===== start =====");
      //   console.log(res);
      //   console.log(res.data);
      //   console.log("===== end =====");
      // })

      // fetch(`http://0.0.0.0:8080/api/users`, {
      //     method: "POST",
      //     headers: {
      //       "Content-Type": "application/json"
      //     },
      //     body: JSON.stringify({
      //       uid: this.state.user.displayName,
      //       name: this.state.user.uid
      //     }),
      // }).then((res) => res.json());



    });
  }

  logout() {
    firebase.auth().signOut();
  }

  render() {
    console.log("===== render start =====");
    console.log(this.state.loading);
    console.log(this.state.user);
    // console.log(firebase.auth().currentUser);
    // console.log(firebase.auth().currentUser.getIdToken(true));
    console.log("===== render end =====");
    if (this.state.loading) return <div>loading</div>;
    return (
      <div>
        Username: {this.state.user && this.state.user.displayName}
        <br />
        {this.state.user ?
          (<button onClick={this.logout}>Logout</button>) :
          (<SignInScreen />)
        }
      </div>
    );
  }
}

export default App;
