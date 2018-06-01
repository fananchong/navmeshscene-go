package NavMeshScene

const HEIGHT_MODE_1 int = 1 // 原始的，精度不是很高，但是没多余消耗。可以通过使Tile Size变小来提高精度
const HEIGHT_MODE_2 int = 2 // 通过公共的StaticScene上，获取精确高度值。
