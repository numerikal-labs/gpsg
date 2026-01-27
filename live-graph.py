import time
import threading
import matplotlib.pyplot as plt
from matplotlib.animation import FuncAnimation
from collections import deque
import pika

WINDOW_SECONDS = 10
MAX_EVENTS = 1000

times = deque(maxlen=MAX_EVENTS)
intervals = deque(maxlen=MAX_EVENTS)
event_lock = threading.Lock()
start_time = None

def rabbitmq_consumer():
    global start_time
    connection = pika.BlockingConnection(
        pika.ConnectionParameters("localhost")
    )
    channel = connection.channel()
    channel.queue_declare(queue="gpsg_queue", durable=True)

    def callback(ch, method, properties, body):
        global start_time
        now = time.time()
        
        with event_lock:
            if start_time is None:
                start_time = now
            
            # Calculate interval from previous event
            if len(times) > 0:
                interval = (now - times[-1]) * 1000  # Convert to milliseconds
                intervals.append(interval)
            else:
                intervals.append(0)  # First event has no interval
            
            times.append(now)

    channel.basic_consume(
        queue="gpsg_queue",
        on_message_callback=callback,
        auto_ack=True,
    )

    channel.start_consuming()

# Start RabbitMQ consumer in background
threading.Thread(target=rabbitmq_consumer, daemon=True).start()

# Setup the plot
fig, ax = plt.subplots(figsize=(12, 6))
fig.suptitle('GPSG Signal Pattern - Waveform', fontsize=16, fontweight='bold')

def update_plot(frame):
    with event_lock:
        if not times or start_time is None or len(times) < 2:
            return
        
        current_times = list(times)
        current_intervals = list(intervals)
        t0 = start_time
    
    now = time.time()
    
    # Clear and setup
    ax.clear()
    
    # Calculate relative times from start
    shifted_times = [t - t0 for t in current_times]
    
    # Create continuous waveform by interpolating between events
    wave_times = []
    wave_values = []
    
    for i in range(len(shifted_times)):
        wave_times.append(shifted_times[i])
        wave_values.append(current_intervals[i])
        
        # Add interpolation point between events (creates the wave connection)
        if i < len(shifted_times) - 1:
            # Add a point just before the next event to create step-like behavior
            next_time = shifted_times[i + 1]
            wave_times.append(next_time - 0.001)  # Just before next event
            wave_values.append(current_intervals[i])  # Hold current value
    
    # Plot the waveform
    ax.plot(wave_times, wave_values, color='blue', linewidth=2, alpha=0.8)
    ax.scatter(shifted_times, current_intervals, c='red', s=50, zorder=5, alpha=0.8, 
               edgecolors='darkred', linewidth=1.5)
    
    # Fill under the curve for better visualization
    ax.fill_between(wave_times, 0, wave_values, alpha=0.2, color='blue')
    
    # Setup axes - show rolling window
    current_relative_time = now - t0
    x_min = max(0, current_relative_time - WINDOW_SECONDS)
    x_max = current_relative_time + 0.5
    
    ax.set_xlim(x_min, x_max)
    
    # Set y-axis limits dynamically based on data
    if current_intervals:
        max_interval = max(current_intervals)
        ax.set_ylim(0, max_interval * 1.2)
    else:
        ax.set_ylim(0, 1000)
    
    ax.set_xlabel('Time (seconds)', fontsize=12)
    ax.set_ylabel('Interval (ms)', fontsize=12)
    ax.grid(True, alpha=0.3)
    ax.axhline(y=0, color='k', linestyle='-', linewidth=0.5)
    
    # Add statistics
    event_count = len(current_times)
    if event_count > 1:
        time_span = current_times[-1] - current_times[0]
        rate = (event_count - 1) / time_span if time_span > 0 else 0
        avg_interval = sum(current_intervals[1:]) / len(current_intervals[1:]) if len(current_intervals) > 1 else 0
        stats_text = f'Events: {event_count} | Rate: {rate:.2f} Hz | Avg Interval: {avg_interval:.1f}ms | Elapsed: {current_relative_time:.1f}s'
    else:
        stats_text = f'Total Events: {event_count} | Waiting for more data...'
    
    ax.text(0.5, 1.05, stats_text, transform=ax.transAxes, 
            ha='center', fontsize=10,
            bbox=dict(boxstyle='round', facecolor='wheat', alpha=0.5))
    
    plt.tight_layout()

# Create animation
ani = FuncAnimation(fig, update_plot, interval=50, cache_frame_data=False)

plt.show()