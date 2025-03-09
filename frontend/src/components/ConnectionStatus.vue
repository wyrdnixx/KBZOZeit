<template>
    <div class="container">

      <table>
        <td> Connection Status</td>
        <td v-if="readyState == 1"  class="green">Socket disconnected {{ readyState }}</td>
        <td v-else class="red">Socket disconnected {{ readyState }}</td>
      </table>

        <!-- <div class="box"><p>1</p></div>
        <div class="box"><p>Connection Status</p></div>
     
        <div v-if="readyState == 1"  class="box green">
           <p> Socket connected {{ readyState }}</p>
            
        </div>
          <div v-else class="box red">            
            <p>Socket disconnected {{ readyState }}</p>
          </div> -->
      
    </div>
    
  </template>
  
  <script>
  export default {
    name: 'ConnectionStatus',
    props: {
        readyState: {
      type: Number,
      required: true,
    }
    },
    computed: {

    },
    watch: {
        readyState(newVal) {
      console.log('Ready state changed to:', newVal);
      // Hier kannst du auf Änderungen reagieren
    }
    },
    data() {
      return {  
        
      };
    },
    mounted() {

    },
    methods: {
  
    }
  }
  </script>
  
  <!-- Add "scoped" attribute to limit CSS to this component only -->
  <style scoped>
    .container {
      display: flex;
      justify-content: center; /* Space between divs can also be 'space-around' or 'space-evenly' */      gap: 10px; /* Adds spacing between the boxes */
    }
    .box {
      width: 20%;
      font-weight: bold;
      border-radius: 5px;
      border: 1px solid #1d8a01; /* Border for the box */

     /*  height: 50px; */
/*       background-color: lightblue; */
    /*   margin: auto; */
   /*    text-align: center; */

  /*     line-height: 100px; */
    /*   font-size: 20px; */
    }

    .red {
      background-color: #ff2234; /* Light red background */
    border: 1px solid #42090f; /* Border for the box */
    color: #192218; 

    }

    .green {
      background-color: #45bb28; /* Light red background */
    border: 1px solid #1d8a01; /* Border for the box */
    color: #192218; 
    }

    table {
      border-collapse: collapse;
      color: #65e057; 

    }
    td {
      border: 1px solid black;
      padding: 8px;
      white-space: nowrap; /* Prevents text from wrapping */
    }
  </style>
  