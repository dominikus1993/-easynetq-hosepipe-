using System.Diagnostics.CodeAnalysis;
using System.IO.IsolatedStorage;

using Akka.Actor;
using Akka.Event;

using EasyNetQ.HosePipe.Types;

using RabbitMQ.Client;
namespace EasyNetQ.HosePipe.Actors;

public sealed record TryRepublish(HosepipeMessage Message)
{
    
}

public sealed class RepublisherActor : ReceiveActor
{
    private IModel _channel;
    private IConnection _connection;
    private readonly IConnectionFactory _connectionFactory;
    
    private readonly ILoggingAdapter _log = Context.GetLogger();
    
    public RepublisherActor(IConnectionFactory factory)
    {
        _connectionFactory = factory;
        _connection = _connectionFactory.CreateConnection();
        _channel = _connection.CreateModel();
        
        OnReady();
    }

    private void OnReady()
    {
        Receive<TryRepublish>(message => {
            _log.Info("Received String message: {0}", message);
            Sender.Tell(message);
        });
    }

    protected override void PreRestart(Exception reason, object message)
    {
        _log.Error(reason, "Actor was stopped {0}", message);
        _connection = _connectionFactory.CreateConnection();
        _channel = _connection.CreateModel();
        base.PreRestart(reason, message);
    }

    protected override void PostStop()
    {
        _channel.Dispose();
        _connection.Dispose();
        base.PostStop();
    }
    
    
    
    
    public static Props Props(IConnectionFactory connectionFactory)
    {
        return Akka.Actor.Props.Create(() => new RepublisherActor(connectionFactory));
    }
}