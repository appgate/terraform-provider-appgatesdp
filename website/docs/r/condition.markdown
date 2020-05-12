---
layout: "appgate"
page_title: "APPGATE: appgate_condition"
sidebar_current: "docs-appgate-resource-condition"
description: |-
   Create a new Condition.
---

# appgate_condition

Create a new Condition..

## Example Usage

```hcl

resource "appgate_condition" "test_condition" {

}

```

## Argument Reference

The following arguments are supported:


* `expression`: (Required) Boolean expression in JavaScript.
* `repeat_schedules`: (Optional) A list of schedules that decides when to reevaluate the Condition. All the scheduled times will be effective. One will not override the other. - It can be a time of the day, e.g. 13:00, 10:25, 2:10 etc. - It can be one of the predefined
  intervals, e.g. 1m, 5m, 15m, 1h. These intervals
  will be always rounded up, i.e. if it&#39;s 15m and the
  time is 12:07 when the Condition is evaluated
  first, then the next evaluation will occur at
  12:15, and the next one will be at
  12:30 and so on.
* `remedy_methods`: (Optional) The remedy methods that will be triggered if the evaluation fails.





## Import

Instances can be imported using the `id`, e.g.

```
$ terraform import appgate_condition d3131f83-10d1-4abc-ac0b-7349538e8300
```
