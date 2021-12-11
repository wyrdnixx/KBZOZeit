<template>
  <div >    
    <p>
      Stempeln
    </p>   
        <button class="btn btn-info" v-on:click="TestBannerMessage()">Test-BannerMessage</button>
  </div>
</template>

<script>

export default {
  name: 'Stempeln',  
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
          this.$parent.showAlert("Stempeln")
      }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
