import {FC} from 'react'

// able to ise all the default body element properties
interface MyCounterProps  extends React.HTMLAttributes<HTMLBodyElement>{
  // custom properties here
  count?: number
}

const MyCounter: FC<MyCounterProps> = ({count, ...props}) => {
  return (
    <h1> {count} </h1>
  )
}
export default MyCounter