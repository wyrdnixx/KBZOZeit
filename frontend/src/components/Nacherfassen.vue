<template>
  <div >    
    <p>
      Nacherfassen
    </p>   
        <button class="btn btn-info" v-on:click="TestBannerMessage()">Test-BannerMessage</button>

        <div>
          <table class="table table-dark">
  <thead>
    <tr>      
      <th scope="col">Datum</th>
      <th scope="col">Von</th>
      <th scope="col">Bis</th>
      <th scope="col">Was</th>
    </tr>
  </thead>
  <tbody>
    <tr>   
      <td><input type="date" id="entryDate" /></td>
      <td><input type="time" id="entryFrom" /></td>
      <td><input type="time" id="entrytimeTo" /></td>      
      <td><input type="text" id="entryText" /></td>      
    </tr>

  </tbody>
</table>
        </div>
  </div>
</template>

<script>

export default {
  name: 'Nacherfassen',  
  props: {
    
  },
  data() {
      return {

      }
  },
  created() {

  },
  methods: {
      TestBannerMessage() {
          this.$parent.showAlert("Nacherfassen")
      }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
