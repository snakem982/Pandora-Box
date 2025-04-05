<template>
  <el-row :gutter="20" class="spark" style="margin-left: 2px;">
    <el-col :span="8">
      <div class="box box1">
        <div class="details">
          <h3>120 Mb/s</h3>
          <h4>{{ $t('home.download') }}</h4>
        </div>
        <apexchart
            id="spark1"
            type="line"
            :options="spark1"
            :series="spark1.series"
            height="100%"
        ></apexchart>
      </div>
    </el-col>
    <el-col :span="8">
      <div class="box box2">
        <div class="details">
          <h3>420 Kb/s</h3>
          <h4>{{ $t('home.upload') }}</h4>
        </div>
        <apexchart
            id="spark2"
            type="line"
            :options="spark2"
            :series="spark2.series"
            height="100%"
        ></apexchart>
      </div>
    </el-col>
    <el-col :span="8">
      <div class="box box3">
        <div class="details">
          <h3>311 Mb</h3>
          <h4>{{ $t('home.memory') }}</h4>
        </div>
        <apexchart
            id="spark3"
            type="line"
            :options="spark3"
            :series="spark3.series"
            height="100%"
        ></apexchart>
      </div>
    </el-col>
  </el-row>
</template>

<script lang="ts" setup>

// setInterval(function () {
//   series[0].data.shift()
//   series[1].data.shift()
//
//   let random1 = Math.floor(Math.random() * 100) + 1;
//   let random2 = Math.floor(Math.random() * 100) + 1;
//
//   series[0].data.push(random1)
//   series[1].data.push(random2)
// }, 2000)

const spark1 = reactive({
  chart: {
    height: 90,
    sparkline: {
      enabled: true
    },
    dropShadow: {
      enabled: true,
      top: 1,
      left: 1,
      blur: 2,
      opacity: 0.5,
    }
  },
  series: [{
    data: [25, 66, 41, 59, 25, 44, 12, 36, 9, 21]
  }],
  stroke: {
    curve: 'smooth'
  },
  markers: {
    size: 0
  },
  grid: {
    padding: {
      top: 20,
      bottom: 10,
      left: 80
    }
  },
  colors: ['#fff'],
  tooltip: {
    theme: 'dark',
    x: {
      show: false
    },
    y: {
      title: {
        formatter: function formatter(val) {
          return '';
        }
      }
    }
  }
})

const spark2 = reactive({
  chart: {
    height: 90,
    sparkline: {
      enabled: true
    },
    dropShadow: {
      enabled: true,
      top: 1,
      left: 1,
      blur: 2,
      opacity: 0.5,
    }
  },
  series: [{
    data: [12, 14, 2, 47, 32, 44, 14, 55, 41, 69]
  }],
  stroke: {
    curve: 'smooth'
  },
  grid: {
    padding: {
      top: 20,
      bottom: 10,
      left: 80
    }
  },
  markers: {
    size: 0
  },
  colors: ['#fff'],
  tooltip: {
    theme: 'dark',
    x: {
      show: false
    },
    y: {
      title: {
        formatter: function formatter(val) {
          return '';
        }
      }
    }
  }
})

const spark3 = reactive({
  chart: {
    height: 90,
    sparkline: {
      enabled: true
    },
    dropShadow: {
      enabled: true,
      top: 1,
      left: 1,
      blur: 2,
      opacity: 0.5,
    }
  },
  series: [{
    data: [47, 45, 74, 32, 56, 31, 44, 33, 45, 19]
  }],
  stroke: {
    curve: 'smooth'
  },
  markers: {
    size: 0
  },
  grid: {
    padding: {
      top: 20,
      bottom: 10,
      left: 80
    }
  },
  colors: ['#fff'],
  xaxis: {
    crosshairs: {
      width: 1
    },
  },
  tooltip: {
    theme: 'dark',
    x: {
      show: false
    },
    y: {
      title: {
        formatter: function formatter(val) {
          return '';
        }
      }
    }
  }
})

</script>

<style scoped>
.spark {
  width: 95%;
}

.spark .box {
  height: 78px;
  padding: 10px 15px;
  text-shadow: 0 1px 1px 1px #666;
  box-shadow: 0 1px 15px 1px rgba(69, 65, 78, 0.08);
  border-radius: 8px;
  position: relative;
}

.spark .box .details {
  position: absolute;
  color: #fff;
  transform: scale(0.8) translate(-15px, -16px);
}

.spark .box1 {
  box-shadow: var(--right-box-shadow);
}

.spark .box2 {
  box-shadow: var(--right-box-shadow);
}

.spark .box3 {
  box-shadow: var(--right-box-shadow);
}
</style>
