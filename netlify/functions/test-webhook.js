exports.handler = async function(event, context) {
  console.log(event)
  console.log(context)
  return {
    statusCode: 200
  }
}