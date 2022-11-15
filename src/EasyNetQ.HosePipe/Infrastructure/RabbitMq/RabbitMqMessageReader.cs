using System.Threading.Channels;

using EasyNetQ.HosePipe.Core.Writers;
using EasyNetQ.HosePipe.Types;

using RabbitMQ.Client;
using RabbitMQ.Client.Events;

namespace EasyNetQ.HosePipe.Infrastructure.RabbitMq;

public sealed class RabbitMqMessageReader : IMessageReader, IDisposable
{
    private readonly IModel _model;

    public RabbitMqMessageReader(IConnection connection)
    {
        _model = connection.CreateModel();
    }

    public IAsyncEnumerable<HosepipeMessage> ReadMessages(Queue queue, CancellationToken cancellationToken)
    {
        var channel = Channel.CreateUnbounded<HosepipeMessage>(new UnboundedChannelOptions() { SingleWriter = true, SingleReader = false});
        var consumer = new AsyncEventingBasicConsumer(_model);
        consumer.Received += async (object sender, BasicDeliverEventArgs ea) =>
        {
            
            await channel.Writer.WriteAsync(null, cancellationToken);
        };
        _model.BasicConsume(queue: queue.Name, autoAck: true, consumer: consumer);
        
        consumer.Unregistered += (_, _) =>
        {
            channel.Writer.TryComplete();
            return Task.CompletedTask;
        };
        consumer.Shutdown += (_, _) =>
        {
            channel.Writer.TryComplete();
            return Task.CompletedTask;
        };
        return channel.Reader.ReadAllAsync(cancellationToken);
    }

    public void Dispose()
    {
        _model.Dispose();
    }
}