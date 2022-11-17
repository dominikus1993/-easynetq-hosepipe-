using System.Threading.Channels;

using EasyNetQ.HosePipe.Core.Writers;
using EasyNetQ.HosePipe.Types;

using OneOf;

using RabbitMQ.Client;
using RabbitMQ.Client.Events;

namespace EasyNetQ.HosePipe.Infrastructure.RabbitMq;

public class RabbitMessageReaderConfiguration
{
    public string? QueueName { get; set; }
}

[GenerateOneOf]
public partial class RabbitMQBasicGetResult : OneOfBase<BasicGetResult, Exception>
{
}

public sealed class RabbitMqMessageReader : IMessageReader, IDisposable
{
    private readonly IModel _model;
    private readonly RabbitMessageReaderConfiguration _readerConfiguration;
    private readonly ILogger<RabbitMqMessageReader> _logger;
    public RabbitMqMessageReader(IConnection connection, RabbitMessageReaderConfiguration readerConfiguration, ILogger<RabbitMqMessageReader> logger)
    {
        _readerConfiguration = readerConfiguration;
        _logger = logger;
        _model = connection.CreateModel();
    }

    public IAsyncEnumerable<HosepipeMessage> ReadMessages(Queue queue, CancellationToken cancellationToken)
    {
        try
        {
            _model.QueueDeclarePassive(_readerConfiguration.QueueName);
        }
        catch (Exception e)
        {
            _logger.LogError(e, "RabbitMq Queue declaration error");
            return AsyncEnumerable.Empty<HosepipeMessage>();
        }
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