//
//  WHOMessageTableViewController.m
//  Hooli
//
//  Created by dylan on 4/19/14.
//  Copyright (c) 2014 whoisdylan. All rights reserved.
//

#import "AFNetworking/AFNetworking.h"
#import "WHOMessageTableViewController.h"
#import "WHONewMessageViewController.h"
#import "WHOMessage.h"
#import "WHOMessageCell.h"

@interface WHOMessageTableViewController () <WHOMessageProtocol>
@property (nonatomic,strong) NSURL* baseURL;
@end

@implementation WHOMessageTableViewController

double kMessageRadius = 1.0;

- (id)initWithStyle:(UITableViewStyle)style WithUserName:(NSString* )userName
{
    self = [super initWithStyle:style];
    if (self) {
        // Custom initialization
        UIColor* hooliColor = [UIColor colorWithRed:109.0/255 green:211.0/255 blue:170.0/255 alpha:1.0];
        UIColor* brownColor = [UIColor colorWithRed:78.0/255 green:46.0/255 blue:40.0/255 alpha:1.0];
        UIColor* goldColor = [UIColor colorWithRed:198.0/255 green:150.0/255 blue:73.0/255 alpha:1.0];
        UIColor* brickColor = [UIColor colorWithRed:207.0/255 green:86.0/255 blue:61.0/255 alpha:1.0];
        [self.tableView setSeparatorInset:UIEdgeInsetsZero];
        [self.tableView setSeparatorColor:hooliColor];
        [self.tableView setBackgroundColor:brownColor];
        UILabel* titleLabel = [[UILabel alloc] init];
        [titleLabel setText:@"Hooli"];
        [titleLabel setFont:[UIFont fontWithName:@"Superclarendon-BlackItalic" size:25.0]];
        [titleLabel setTextColor:goldColor];
//        [titleLabel setAlpha:0.75];
        [titleLabel.layer setShadowColor:[UIColor darkGrayColor].CGColor];
        [titleLabel.layer setShadowOffset:(CGSize) { .width = 1.0, .height = 1.0 }];
        [titleLabel.layer setShadowRadius:1.0];
        [titleLabel.layer setShadowOpacity:.65];
        [titleLabel sizeToFit];
        [self.navigationItem setTitleView:titleLabel];
        self.messages = [NSMutableArray array];
        self.userName = userName;
        
        self.locationManager = [[CLLocationManager alloc] init];
        self.locationManager.desiredAccuracy = kCLLocationAccuracyBest;
        self.locationManager.delegate = self;
        [self.locationManager startUpdatingLocation];
    }
    return self;
}

- (void)viewDidLoad
{
    [super viewDidLoad];
    
    // Uncomment the following line to preserve selection between presentations.
    // self.clearsSelectionOnViewWillAppear = NO;
    
    // Uncomment the following line to display an Edit button in the navigation bar for this view controller.
    // self.navigationItem.rightBarButtonItem = self.editButtonItem;
    
//    UIColor* hooliColor =[UIColor colorWithRed:109.0/255 green:211.0/255 blue:170.0/255 alpha:1.0];
    UIColor* brownColor = [UIColor colorWithRed:78.0/255 green:46.0/255 blue:40.0/255 alpha:1.0];
    UIBarButtonItem* newMessageButton = [[UIBarButtonItem alloc] initWithBarButtonSystemItem:UIBarButtonSystemItemCompose target:self action:@selector(newMessage:)];
    [newMessageButton setTintColor:brownColor];
    self.navigationItem.rightBarButtonItem = newMessageButton;
    
    [self.tableView setRowHeight:150.0];
    [self.tableView registerNib:[UINib nibWithNibName:@"WHOMessageCell" bundle:[NSBundle mainBundle]] forCellReuseIdentifier:@"MessageCell"];
    
    self.refreshControl = [[UIRefreshControl alloc] init];
    [self.refreshControl addTarget:self action:@selector(updateMessageList) forControlEvents:UIControlEventValueChanged];
    
    self.baseURL = [NSURL URLWithString:@"http://192.168.1.19:9009/proposer/"];
    
    [self updateMessageList];
//    [NSTimer scheduledTimerWithTimeInterval:10.0 target:self selector:@selector(updateMessageList) userInfo:nil repeats:YES];
}

- (void)updateMessageList {
    NSLog(@"Updating message list");
    
    //using REST
//    NSDictionary* parameters = @{@"MessageText" : message, @"Author" : self.userName, @"Latitude" : [NSNumber numberWithDouble: self.userLocation.coordinate.latitude], @"Longitude" : [NSNumber numberWithDouble: self.userLocation.coordinate.longitude]};
    
    AFHTTPSessionManager *manager = [[AFHTTPSessionManager alloc] initWithBaseURL:self.baseURL];
    manager.responseSerializer = [AFJSONResponseSerializer serializer];
//    manager.responseSerializer.acceptableContentTypes = [NSSet setWithObject:@"application/json"];
//    manager.responseSerializer = [AFHTTPResponseSerializer serializer];
//    manager.responseSerializer.acceptableContentTypes = [NSSet setWithObject:@"text/plain"];
//    manager.requestSerializer = [AFJSONRequestSerializer serializer];
    manager.requestSerializer = [AFHTTPRequestSerializer serializer];
//    manager.requestSerializer = [AFJSONRequestSerializer serializer];
//    [manager.requestSerializer setValue:@"text/html" forHTTPHeaderField:@"Content-Type"];
    [manager.requestSerializer setValue:@"application/json" forHTTPHeaderField:@"Accept"];
    
    [manager GET:@"messages" parameters:NULL success:^(NSURLSessionDataTask *task, id responseObject) {
        NSLog(@"REST GET success!");
        NSLog(@"received response containing: %@", (NSArray*) responseObject);
        for (NSDictionary* element in responseObject) {
//            NSLog(@"fieldtest:%@",element[@"Author"]);
            NSNumber* lat = element[@"Latitude"];
            NSNumber* lon = element[@"Longitude"];
            CLLocation* loc = [[CLLocation alloc] initWithLatitude:lat.doubleValue longitude:lon.doubleValue];
            if ([self isMessageWithinRangeAtLocation:loc]) {
                NSLog(@"Message is within region");
                WHOMessage* mess = [[WHOMessage alloc] initWithMessage:element[@"MessageText"] Author:element[@"Author"] Location:loc];
                if (![self.messages containsObject:mess]) {
                    NSLog(@"Message is not already in list");
                    [self.messages addObject:mess];
                }
                /*BOOL messageExists = NO;
                for (WHOMessage* prevMess in self.messages) {
                    if ([prevMess isEqualToObject:mess]) {
                        NSLog(@"message exists!");
                        messageExists = YES;
                        break;
                    }
                }
                if (!messageExists) {
                    NSLog(@"message doesn't exist!");
                    [self.messages addObject:mess];
                }*/
            }
        }
        [self.refreshControl endRefreshing];
        [self.tableView reloadData];
    } failure:^(NSURLSessionDataTask *task, NSError *error) {
        NSLog(@"REST GET failure wih error: %@", error);
        [self.refreshControl endRefreshing];
    }];
//    [self.refreshControl endRefreshing];
}

- (void)locationManager:(CLLocationManager *)manager didUpdateLocations:(NSArray *)locations {
    CLLocation *currLoc = [locations firstObject];
//    NSLog(@"User location updated: %@", currLoc);
    self.userLocation = currLoc;
    
}

- (void)newMessage:(id) sender {
    WHONewMessageViewController* form = [[WHONewMessageViewController alloc] init];
    form.delegate = self;
    [UIView animateWithDuration:0.7
                     animations:^{
                         [UIView setAnimationCurve:UIViewAnimationCurveEaseInOut];
                         [UIView setAnimationTransition:UIViewAnimationTransitionCurlUp forView:self.navigationController.view cache:NO];
                     }];
    [self.navigationController pushViewController:form animated:NO];
    
}

- (void)receivedNewMessage:(NSString *)message {
//    NSLog(@"received new message %@",message);
//    WHOMessage* messageObj = [[WHOMessage alloc] initWithMessage:message Author:self.userName Location:self.userLocation];
//    NSString* latitude = [[NSString alloc] initWithFormat:@"%f", self.userLocation.coordinate.latitude];
//    NSString* longitude = [[NSString alloc] initWithFormat:@"%f", self.userLocation.coordinate.longitude];
    
    //using REST
    NSDictionary* parameters = @{@"MessageText" : message, @"Author" : self.userName, @"Latitude" : [NSNumber numberWithDouble: self.userLocation.coordinate.latitude], @"Longitude" : [NSNumber numberWithDouble: self.userLocation.coordinate.longitude]};
    
    AFHTTPSessionManager *manager = [[AFHTTPSessionManager alloc] initWithBaseURL:self.baseURL];
//    manager.responseSerializer = [AFJSONResponseSerializer serializer];
    manager.responseSerializer = [AFHTTPResponseSerializer serializer];
//    manager.responseSerializer.acceptableContentTypes = [NSSet setWithObject:@"text/plain"];
    manager.requestSerializer = [AFJSONRequestSerializer serializer];
    
    [manager POST:@"messages" parameters:parameters success:^(NSURLSessionDataTask *task, id responseObject) {
        NSLog(@"REST POST success!");
    } failure:^(NSURLSessionDataTask *task, NSError *error) {
        NSLog(@"REST POST failure wih error: %@", error);
    }];
}

- (BOOL)isMessageWithinRangeAtLocation:(CLLocation*)location {
    CLLocationDistance distanceInMeters = fabs([location distanceFromLocation:self.userLocation]);
    double distanceInMiles = distanceInMeters/1609.344;
    if (distanceInMiles <= kMessageRadius) {
        return YES;
    }
    else {
        return NO;
    }
}

- (NSString *)distanceBetweenUserAndLocation:(CLLocation *)location {
    CLLocationDistance CLDistance = fabs([location distanceFromLocation:self.userLocation]);
    NSString* distance = [NSString stringWithFormat:@"%.1f miles away",(CLDistance/1609.344)];
    return distance;
}

- (void)didReceiveMemoryWarning
{
    [super didReceiveMemoryWarning];
    // Dispose of any resources that can be recreated.
}

#pragma mark - Table view data source

- (NSInteger)numberOfSectionsInTableView:(UITableView *)tableView
{
    // Return the number of sections.
    return 1;
}


- (NSInteger)tableView:(UITableView *)tableView numberOfRowsInSection:(NSInteger)section
{
    // Return the number of rows in the section.
    return self.messages.count;
}

- (UITableViewCell *)tableView:(UITableView *)tableView cellForRowAtIndexPath:(NSIndexPath *)indexPath
{
    WHOMessageCell *cell = [tableView dequeueReusableCellWithIdentifier:@"MessageCell" forIndexPath:indexPath];
    WHOMessage* message = [self.messages objectAtIndex:indexPath.row];
    cell.messageLabel.text = message.message;
    cell.authorLabel.text = message.author;
    cell.distanceLabel.text = [self distanceBetweenUserAndLocation:message.location];
    return cell;
}

- (CGFloat)tableView:(UITableView *)tableView heightForRowAtIndexPath:(NSIndexPath *)indexPath {
    return 150.0;
}

/*
// Override to support conditional editing of the table view.
- (BOOL)tableView:(UITableView *)tableView canEditRowAtIndexPath:(NSIndexPath *)indexPath
{
    // Return NO if you do not want the specified item to be editable.
    return YES;
}
*/

/*
// Override to support editing the table view.
- (void)tableView:(UITableView *)tableView commitEditingStyle:(UITableViewCellEditingStyle)editingStyle forRowAtIndexPath:(NSIndexPath *)indexPath
{
    if (editingStyle == UITableViewCellEditingStyleDelete) {
        // Delete the row from the data source
        [tableView deleteRowsAtIndexPaths:@[indexPath] withRowAnimation:UITableViewRowAnimationFade];
    } else if (editingStyle == UITableViewCellEditingStyleInsert) {
        // Create a new instance of the appropriate class, insert it into the array, and add a new row to the table view
    }   
}
*/

/*
// Override to support rearranging the table view.
- (void)tableView:(UITableView *)tableView moveRowAtIndexPath:(NSIndexPath *)fromIndexPath toIndexPath:(NSIndexPath *)toIndexPath
{
}
*/

/*
// Override to support conditional rearranging of the table view.
- (BOOL)tableView:(UITableView *)tableView canMoveRowAtIndexPath:(NSIndexPath *)indexPath
{
    // Return NO if you do not want the item to be re-orderable.
    return YES;
}
*/

/*
#pragma mark - Navigation

// In a storyboard-based application, you will often want to do a little preparation before navigation
- (void)prepareForSegue:(UIStoryboardSegue *)segue sender:(id)sender
{
    // Get the new view controller using [segue destinationViewController].
    // Pass the selected object to the new view controller.
}
*/

@end
