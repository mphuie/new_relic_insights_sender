# new_relic_insights_sender
Simple CLI tool for sending custom events to New Relic Insights

Pipe output of command(s) or just echo text to create an event in Insights.

Run it on a schedule and create a dashboard and/or alert in New Relic.


## Example usage

Count the number of errors and send an event to New Relic Insights

    grep error yourapp.log | wc -l | ./nr_insights_sender   # by default uses the `value` key
    
    grep error yourapp.log | wc -l | ./nr_insights_sender errorCount  # specify a your own key to send the data as
    
    
    
Will appear in your new relic insights as

EventType: <from config>

    {
        value: <value>   or errorCount: <value>
        timestamp: <auto generated by New relic>
        
        ... any additional k/v from your config
    }
 

## Configuration

Create a `new_relic_insights_sender.yaml` config file with your New Relic account ID and Insights insert key.
Create your default values to be sent in the `EventValues` key.  You will at least need an `eventType`


Sample config 

    NewRelicAccountId: <your account id>
    NewRelicInsertKey: <your insert key>
    EventValues:
      eventType: myapp                       # required
      description: error count in logs       # all other fields optional
      host: exampleHost      
