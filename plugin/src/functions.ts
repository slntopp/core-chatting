import { NotificationApi } from "naive-ui"

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

export function getRelativeTime(timestamp: number, now: number) {
  const timeDifference = (now - timestamp) / 1000;
  const minutesDifference = Math.floor(timeDifference / 60);

  if (minutesDifference >= 4320) {
    return new Date(timestamp).toLocaleDateString();
  } else if (minutesDifference >= 1440) {
    const daysDifference = Math.floor(minutesDifference / 1440);
    return `${daysDifference} days ago`;
  } else if (minutesDifference >= 60) {
    const hoursDifference = Math.floor(minutesDifference / 60);
    return `${hoursDifference} hours ago`;
  } else if (minutesDifference > 0) {
    return `${minutesDifference} minutes ago`;
  } else {
    return 'just now';
  }
}