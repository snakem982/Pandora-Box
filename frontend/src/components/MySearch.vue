<script setup lang="ts" xmlns="http://www.w3.org/1999/html">
import {ArrowLeft, ArrowRight, Close, Search} from "@element-plus/icons-vue";
</script>

<template>
  <div class="search-container">
    <span class="back">
      <el-icon><ArrowLeft/></el-icon>
    </span>
    <span class="forward">
      <el-icon><ArrowRight/></el-icon>
    </span>
    <input type="text"
           id="search-input"
           autocapitalize="off"
           autocomplete="off"
           autocorrect="off"
           spellcheck="false"
           placeholder="搜索节点">

    <span class="search">
      <el-icon>
        <Search/>
      </el-icon>
    </span>
    <span class="clear"
          id="clear-button">
      <el-icon>
        <Close/>
      </el-icon>
    </span>

    <span class="minus">
      <el-icon>
        <icon-mdi-card-minus-outline/>
      </el-icon>
    </span>


    <div class="dropdown" id="dropdown">
      <ul>
        <li class="group">A</li>
        <li>Alice</li>
        <li>Adam</li>
        <li class="group">B</li>
        <li>Bob</li>
        <li>Betty</li>
        <li class="group">C</li>
        <li>Charlie</li>
        <li>Chloe</li>
      </ul>
    </div>
  </div>

</template>

<style scoped>
.search-container {
  left: 230px;
  padding-top: 25px;
  z-index: 100;
  border-top-left-radius: 10px;
  border-top-right-radius: 10px;
}

.search-container input {
  width: 138px;
  padding: 10px 30px 10px 40px;
  border: none;
  border-radius: 10px;
  background-color: rgba(255, 255, 255, 0.06);
  color: var(--text-color);
  font-size: 12px;
  margin-left: 26px;
  margin-top: -8px;
  position: absolute;
}

.search-container input:focus {
  outline: none;
}

.search-container .back {
  margin-left: 35px;
  font-size: 18px;
  color: var(--text-color);
}

.search-container .forward {
  margin-left: 25px;
  font-size: 18px;
  color: var(--text-color);
}

.search-container .search {
  margin-left: 38px;
  margin-top: 28px;
  font-size: 18px;
  color: var(--text-color);
}

.search-container .clear {
  margin-left: 308px;
  margin-top: -18px;
  font-size: 12px;
  color: var(--text-color);
  cursor: pointer;
  display: block;
  position: absolute;
}

.search-container .minus {
  margin-right: 26px;
  float: right;
  font-size: 16px;
  color: var(--text-color);
}

.dropdown {
  top: 50px;
  left: 0;
  width: 100%;
  border-radius: 5px;
  background-color: #fff;
  display: none;
  z-index: 1000;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.2);
}

.dropdown ul {
  margin: 0;
  padding: 0;
  list-style: none;
}

.dropdown ul .group {
  padding: 5px;
  font-weight: bold;
  color: #333;
  background-color: #f9f9f9;
  border-bottom: 1px solid #eee;
}

.dropdown ul li {
  padding: 10px;
  cursor: pointer;
}

.dropdown ul li:hover {
  background-color: #f0f0f0;
}
</style>