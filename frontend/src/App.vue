<template>
  <div class="cBody"
       :style="{'--body-background-image': backgroundImage}"
       key="pandora-box-body"
  >
    <div class="left">
      <div class="top-title">
        <div class="top-icon"></div>
        <span class="top-title-text">Pandora-Box</span>
      </div>

      <MyNav></MyNav>
      <MyRule></MyRule>
      <MyProxy></MyProxy>
      <MySecNav></MySecNav>
      <MyBottom></MyBottom>

    </div>

    <div class="right">
      <router-view/>
      <MyDrop></MyDrop>
    </div>
  </div>
</template>

<script setup lang="ts">
import MyNav from "@/components/menu/MyNav.vue";
import MySecNav from "@/components/menu/MySecNav.vue";
import MyDrop from "@/components/MyDrop.vue";
import MyBottom from "@/components/menu/MyBottom.vue";
import {useWebStore} from "@/store/webStore";

const backgroundImage = ref("url('/images/quang.jpg')")
// backgroundImage.value = "radial-gradient(circle, #2c3e50, #4a5568)"
// backgroundImage.value = "linear-gradient(to bottom, #1b2735, #3b4b69)"
// backgroundImage.value = "linear-gradient(to bottom, #3a3d40, #181818)"
backgroundImage.value = "linear-gradient(to bottom, #434343, #000000)"


const webStore = useWebStore()

onMounted(()=>{
  console.log("=========",webStore.baseUrl)
  console.log(webStore.secret)
})


</script>


<style>
.cBody {
  margin: 0;
  display: flex;
  height: 100vh;
  background-image: var(--body-background-image);
  background-size: cover;
  background-position: center;
  color: var(--text-color);

}

.left {
  padding-right: 18px;
}

.right {
  overflow-y: hidden;
  overflow-x: hidden;
  position: relative;
  width: 100%;
  flex-grow: 1;
  margin: 15px 15px 15px 0;
  border-radius: 10px;
  background-color: var(--right-bg-color);
  color: var(--text-color);
}

.top-title {
  margin-top: 40px;
  margin-left: 20px;
}

.top-icon {
  width: 28px;
  height: 28px;
  background-image: url("@/assets/images/appicon.png");
  background-size: cover;
  background-position: center;
}

.top-title-text {
  position: absolute;
  margin-left: 40px;
  margin-top: -25px;
}

</style>