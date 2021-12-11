<template>
  <div>

    <h1>{{ msg }}</h1>
    <p>
      Main App
    </p>

    <!-- Toast Message banner  -->

    <b-alert
      :show="AlertDismissCountDown"
      dismissible
      variant="warning"
      @dismissed="AlertDismissCountDown=0"
      @dismiss-count-down="countDownChanged"
    >
      <p> {{this.AlertWebsiteAlertMessage}} </p>
        <!-- {{ AlertDismissCountDown }} </p> -->
      <b-progress
        variant="warning"
        :max="AlertDismissSecs"
        :value="AlertDismissCountDown"
        height="4px"
      ></b-progress>
    </b-alert>

<!-- Toast Message banner  -->

    <button class="btn btn-info" v-on:click="TestChangeAuth()">Test-ChangeAuthenticated</button>

    <div v-if="!this.UserAuthenticated">
      <Register />
    </div>
    <div v-else>
      <TimeAccounting />
    </div>
  </div>
</template>

<script>
import Register from './Register.vue'
import TimeAccounting from './TimeAccounting.vue'


export default {
  name: 'Main',
  components: {
    Register,
    TimeAccounting
  },
  props: {
    msg: String
  },
  data() {
    return {
      UserAuthenticated: "",
      AlertDismissSecs: 4,
      AlertDismissCountDown: 0,
      AlertShowDismissibleAlert: false,
      AlertWebsiteAlertMessage: ""   
    }
    
  },
  created() {

    this.UserAuthenticated = true
  },
  methods: {
    TestChangeAuth() {
      console.log("TestChangeAuth")
      if (this.UserAuthenticated) {
        console.log("setting false")
        this.UserAuthenticated = false
      } else {
        console.log("setting true")
        this.UserAuthenticated = true
      }
    },
    countDownChanged(AlertDismissCountDown) {
      this.AlertDismissCountDown = AlertDismissCountDown
    },
    showAlert(msg) {

      this.AlertWebsiteAlertMessage = msg
      this.AlertDismissCountDown = this.AlertDismissSecs
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
</style>
