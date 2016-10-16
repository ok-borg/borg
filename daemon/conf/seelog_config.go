package conf

const Seelog = `
<seelog>
  <outputs>
    <console formatid="colored"/>
  </outputs>
  <formats>
    <format id="colored"  format="%Date(2006 Jan 02/3:04:05.00 PM MST) (%File) [%EscM(36)%LEVEL%EscM(39)] %Msg%n%EscM(0)"/>
  </formats>
</seelog>
`
