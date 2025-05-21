<template>
  <div class="cBody"
       :style="{ backgroundImage: currentBackground }"
       key="pandora-box-body"
  >
    <div class="left">
      <div :class="isWindows?'top-title win':'top-title'">
        <div class="top-icon"></div>
        <span class="top-title-text">Pandora-Box</span>
      </div>
      <MyEvent/>
      <MyNav/>
      <MyRule/>
      <MyProxy/>
      <MySecNav/>
      <MyBottom/>
    </div>

    <div class="right">
      <router-view/>
      <MyDrop/>
    </div>
  </div>
</template>


<script setup lang="ts">
import {onMounted, ref} from 'vue';
import {useMenuStore} from "@/store/menuStore";

const menuStore = useMenuStore();

// 当前背景
const defaultBackground = "linear-gradient(to bottom, #434343, #000000)"
const currentBackground = ref(defaultBackground);

// 预加载背景
function preloadBackgroundImage(bg: string) {
  if (bg.startsWith('url(')) {
    const imgUrl = bg.slice(5, -2); // 提取 url('/path') 中的地址
    const img = new Image();
    img.src = imgUrl;
    img.onload = () => {
      currentBackground.value = bg;
    };
    img.onerror = () => {
      console.error(`Failed to load background image: ${imgUrl}`);
      currentBackground.value = defaultBackground;
    };
  } else {
    currentBackground.value = bg; // 直接应用渐变背景
  }
}

const isWindows = ref(false)
onMounted(() => {
  preloadBackgroundImage(menuStore.background);
  // @ts-ignore
  if (window["pxShowBar"]) {
    isWindows.value = true;
  }
});

watch(() => menuStore.background, (nextBackground) => {
  preloadBackgroundImage(nextBackground);
});

</script>


<style>
.cBody {
  margin: 0;
  display: flex;
  height: 100vh;
  color: var(--text-color);
  background-color: #000;
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  background-attachment: fixed;
  transition: background-image 0.5s ease-in-out, background-color 0.5s ease-in-out;
  background-color: rgba(0, 0, 0, 0.2);
  background-blend-mode: multiply;
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
  padding-top: 40px;
  padding-left: 23px;
  -webkit-app-region: drag;
}

.win {
  padding-top: 35px;
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
  margin-top: -22px;
}
</style>
