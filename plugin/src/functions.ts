import { NotificationApi } from 'naive-ui'
import { Status } from './connect/cc/cc_pb'

export function addToClipboard (text: string, notification?: NotificationApi) {
  if (navigator?.clipboard) {
    navigator.clipboard.writeText(text)
      .then(() => {
        if (!notification) {
          alert('Text copied')
          return
        }
        notification.success({ content: 'Text copied', duration: 3000 })
      })
      .catch((res) => {
        console.error(res)
      })
  } else {
    if (!notification) {
      alert('Clipboard is not supported')
      return
    }
    notification.error({ content: 'Clipboard is not supported' })
  }
}

export function getImageUrl(name: string) {
  const icons = import.meta.glob(`/**/*.png`)
  let key = Object.keys(icons).find((key) => key.includes(`/${name}`)) ?? ''

  if (!key.includes('cc.ui')) key = `/cc.ui${key}`
  return key.replace('/dist', '').replace('/public', '')
}

export function getRelativeTime(timestamp: number, now: number, isLifetime?: boolean) {
  const timeDifference = (now - timestamp) / 1000;
  const minutesDifference = Math.floor(timeDifference / 60);

  if (minutesDifference >= 4320 && !isLifetime) {
    return new Date(timestamp).toLocaleDateString();
  } else if (minutesDifference >= 1440) {
    const daysDifference = Math.floor(minutesDifference / 1440);

    return `${daysDifference} days${(isLifetime) ? '' : ' ago'}`;
  } else if (minutesDifference >= 60) {
    const hoursDifference = Math.floor(minutesDifference / 60);

    return `${hoursDifference} hours${(isLifetime) ? '' : ' ago'}`;
  } else if (minutesDifference > 0) {
    return `${minutesDifference} minutes${(isLifetime) ? '' : ' ago'}`;
  } else {
    return 'just now';
  }
}

export function getStatusColor (status: Status) {
  switch (status) {
    case Status.NEW:
      return '#5084ff'
    case Status.ON_HOLD:
      return '#00dbff'
    case Status.OPEN:
      return '#1ea01e'
    case Status.IN_PROGRESS:
      return '#00ffaa'
    case Status.CUSTOMER_REPLY:
      return '#ffcc55'
    case Status.RESOLVE:
    case Status.ANSWERED:
      return '#ff8300'
    case Status.CLOSE:
      return '#e23535'
    default:
      return undefined
  }
}
