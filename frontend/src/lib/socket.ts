let scheme = 'ws'
const location = document.location

if (location.protocol === 'https:') {
  scheme += 's'
}

const serverUrl = `${scheme}://${location.hostname}:8080`
export const websocket = new WebSocket(serverUrl)
// export const websocket = new WebSocket('ws://localhost:8080/')
// function closeWebsocket() {
//   websocket.removeEventListener('open', onConnect)
//   websocket.removeEventListener('close', reConnect)
//   websocket.removeEventListener('message', handleMessage)
// }
// const connectWebSocket = () => {
//   websocket.addEventListener('open', onConnect)
//   websocket.addEventListener('close', reConnect)
//   websocket.addEventListener('message', handleMessage)
// }
// const reConnect = () => {
//   onDisconnect()
//   setTimeout(() => {
//     connectWebSocket()
//   }, 3000)
// }
// function onConnect() {
//   console.log('connected')
//   //   WSState.IsConnected = true
// }

// function onDisconnect() {
//   //   WSState.IsConnected = false
// }

// function handleMessage(this: WebSocket, ev: MessageEvent<any>) {
//   throw new Error('Function not implemented.')
// }
// const playerOptions = () => {
//   const options = []
//   for (let index = 1; index <= 9; index++) {
//     options.push({ id: index, number: index.toString() })
//   }
//   const a = options.filter((item) => item.id % 2 === 0)
//   const b = options.filter((item) => item.id % 2 !== 0)
//   // console.log(a);
//   return { a, b }
// }
// export function calculateWinner(squares: Detail[]) {
//   const lines = [
//     [0, 1, 2],
//     [3, 4, 5],
//     [6, 7, 8],
//     [0, 3, 6],
//     [1, 4, 7],
//     [2, 5, 8],
//     [0, 4, 8],
//     [2, 4, 6]
//   ]
//   for (let i = 0; i < lines.length; i++) {
//     const [a, b, c] = lines[i]
//     // console.log(Number(squares[a].number) + Number(squares[b].number) + Number(squares[c].number));
//     // console.log(Number(squares[a].number), Number(squares[b].number), Number(squares[c].number));

//     if (Number(squares[a].number) + Number(squares[b].number) + Number(squares[c].number) === 15) {
//       console.log(Number(squares[a].number), Number(squares[b].number), Number(squares[c].number))
//       return true
//     }
//   }
//   return null
// }
