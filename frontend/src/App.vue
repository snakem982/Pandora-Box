<template>
  <div class="cBody" :style="{'--body-background-image': backgroundImage}">
    <div class="left">
      <div class="top-title">
        <div class="top-icon"></div>
        <span class="top-title-text">Pandora-Box</span>
      </div>

      <MyNav></MyNav>
      <MyRule></MyRule>
      <MyProxy></MyProxy>
      <MySecNav></MySecNav>

      <div class="bottom-text">
        <el-icon>
          <icon-mdi-power/>
        </el-icon>
        <el-icon style="margin-left:20px">
          <icon-mdi-translate/>
        </el-icon>
        <el-icon style="margin-left:20px">
          <icon-mdi-tshirt-crew-outline/>
        </el-icon>
      </div>
    </div>

    <div class="right">
      <router-view />
      <MyDrop></MyDrop>
    </div>
  </div>
</template>

<script setup lang="ts">
import MyNav from "@/components/menu/MyNav.vue";
import MySecNav from "@/components/menu/MySecNav.vue";
import MyDrop from "@/components/MyDrop.vue";

const backgroundImage = ref("url('/src/assets/images/sea.jpg')")
// backgroundImage.value = "radial-gradient(circle, #2c3e50, #4a5568)"
// backgroundImage.value = "linear-gradient(to bottom, #1b2735, #3b4b69)"


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


.bottom-text {
  position: absolute;
  bottom: 18px;
  left: 22px;
  color: var(--text-color);
  font-size: 20px;
}

</style>